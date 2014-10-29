package controllers

import (
	"df/assets/models"
	. "df/assets/utils"
	"df/errmap"
	img "github.com/3d0c/imgproc"
	"github.com/3d0c/oid"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/encoder"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type File struct{}

var guid *regexp.Regexp = regexp.MustCompile("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")

func (*File) Construct(args ...interface{}) interface{} {
	return &File{}
}

func (this *File) Find(enc encoder.Encoder, w http.ResponseWriter, r *http.Request, p martini.Params, options models.URLOptionsScheme) (int, []byte) {
	if guid.MatchString(p["id"]) || options.Ids != "" {
		return this.ServeGeometry(enc, p, options, w, r)
	}

	if NewTypeFromFlag(oid.ObjectIdHex(p["id"]).Flag()).IsImage {
		return this.ServeImage(enc, p, options, w)
	}

	name := AppConfig.StorageFilePath(p["id"])
	http.ServeFile(w, r, name)

	return 200, []byte{}
}

func (this *File) ServeGeometry(enc encoder.Encoder, params martini.Params, options models.URLOptionsScheme, w http.ResponseWriter, r *http.Request) (int, []byte) {
	model := (*models.Geometry).Construct(nil).(*models.Geometry)

	ids := []string{}

	if options.Ids != "" {
		for _, id := range strings.Split(options.Ids, ",") {
			if guid.MatchString(id) {
				ids = append(ids, id)
			}
		}
	} else {
		ids = append(ids, params["id"])
	}

	if len(ids) > 1 && len(options.ToMap()) == 0 {
		options.Height = 175
	}

	result := model.FindAll(ids, options)

	if strings.Contains("application/json", r.Header.Get("Accept")) {
		return http.StatusOK, encoder.Must(enc.Encode(result))
	}

	w.Header().Set("Content-Type", "application/octet-stream")

	if result == nil || len(result) == 0 {
		return http.StatusNotFound, []byte{}
	}

	pmap := options.ToMap()
	// name := AppConfig.StorageFilePath(result.Base.Id) + options.ToHash()

	names := []string{}
	for _, geometry := range result {
		var prefix string
		if len(ids) > 1 {
			prefix = options.Ids + "."
		}

		name := AppConfig.StoragePath(geometry.Base.Id) + prefix + geometry.Base.Id + options.ToHash()

		names = append(names, name)
	}

	log.Println("Serving files:", names)

	// check only for first file in set, because if there is no such file, whole set
	// havn't been ever morphed
	if _, err := os.Stat(names[0]); err != nil {
		models.Morph(names, result, pmap, options)
	}

	var sizes []string

	for _, name := range names {
		fi, err := os.Stat(name)
		if err != nil {
			log.Println("File not found after morphing, error:", err)
			return http.StatusNotFound, []byte{}
		}

		sizes = append(sizes, strconv.FormatInt(fi.Size(), 10))
	}

	w.Header().Set("Df-Sizes", strings.Join(sizes, ","))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(200)

	for _, name := range names {
		file, err := os.Open(name)
		if err != nil {
			log.Println("File not found after morphing, error:", err)
			return http.StatusNotFound, []byte{}
		}

		io.Copy(w, file)
		file.Close()
	}

	return http.StatusOK, []byte{}
}

func (this *File) ServeImage(enc encoder.Encoder, params martini.Params, options models.URLOptionsScheme, w http.ResponseWriter) (int, []byte) {
	options = options.SetDefaults()

	base := img.NewSource(AppConfig.StorageFilePath(params["id"]))
	if base == nil {
		return http.StatusNotFound, []byte{}
	}

	target := &img.Options{
		Base:    base,
		Scale:   img.NewScale(options.Scale),
		Method:  3,
		Quality: options.Quality,
	}

	var typ Type
	if options.Format == "" {
		typ = NewTypeFromFlag(oid.ObjectIdHex(params["id"]).Flag())
	} else {
		typ = NewTypeFromName(options.Format)
	}

	target.Format = typ.Format()

	w.Header().Set("Content-Type", typ.ContentType())

	return http.StatusOK, img.Proc(target)
}

func (this *File) Create(enc encoder.Encoder, req *http.Request) (int, []byte) {
	if !strings.Contains(req.Header.Get("Content-Type"), "multipart/form-data") {
		return http.StatusBadRequest, encoder.Must(enc.Encode(errmap.Err("Expected ContentType is multipart/form-data.")))
	}

	result := make([]models.FileScheme, 0)

	if err := req.ParseMultipartForm(64 << 20); err != nil {
		log.Println("ParseMultipartForm error:", err)
		return http.StatusInternalServerError, []byte{}
	}

	isErr := func(err error, name string) bool {
		if err != nil {
			log.Println("Unable to process uploaded file:", name, ",error:", err)
			result = append(result, models.FileScheme{Name: name, Err: "Unable to porcess file. Try again."})
			return true
		}

		return false
	}

	for _, fileHeaders := range req.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			file, err := fileHeader.Open()
			if isErr(err, fileHeader.Filename) {
				continue
			}

			typ := NewTypeFromExt(fileHeader.Filename)

			id := oid.NewObjectId(typ.Flag(), AppConfig.NodeId()).String()
			log.Printf("%s, node: %d, flag: %d\n", id, AppConfig.NodeId(), oid.ObjectIdHex(id).Flag())

			buf, err := ioutil.ReadAll(file)
			if isErr(err, fileHeader.Filename) {
				continue
			}

			err = ioutil.WriteFile(AppConfig.StorageFor(id), buf, os.ModePerm)
			if isErr(err, fileHeader.Filename) {
				continue
			}

			result = append(result, models.FileScheme{Id: id, Name: fileHeader.Filename})
		}
	}

	return 200, encoder.Must(enc.Encode(result))
}
