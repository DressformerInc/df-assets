package controllers

import (
	"df/assets/models"
	. "df/assets/utils"
	"df/errmap"
	img "github.com/3d0c/imgproc"
	"github.com/3d0c/oid"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/encoder"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type File struct{}

func (*File) Construct(args ...interface{}) interface{} {
	return &File{}
}

func (this *File) Find(enc encoder.Encoder, w http.ResponseWriter, r *http.Request, p martini.Params, options models.URLOptionsScheme) (int, []byte) {
	guid := regexp.MustCompile("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")
	if guid.MatchString(p["id"]) {
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

	result := model.Find(params["id"])

	if strings.Contains("application/json", r.Header.Get("Accept")) {
		return http.StatusOK, encoder.Must(enc.Encode(result))
	}

	w.Header().Set("Content-Type", "application/octet-stream")

	if result == nil || result.Base.Id == "" {
		return http.StatusNotFound, []byte{}
	}

	pmap := options.ToMap()
	name := AppConfig.StorageFilePath(result.Base.Id) + options.ToHash()
	log.Println("Serving file:", name)

	if _, err := os.Stat(name); err != nil {
		result.Morph(name, pmap, options)
	}

	http.ServeFile(w, r, name)

	return http.StatusNotFound, []byte{}
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
