package controllers

import (
	"df/assets/models"
	. "df/assets/utils"
	"df/errmap"
	"github.com/3d0c/oid"
	"github.com/martini-contrib/encoder"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type File struct{}

func (*File) Construct(args ...interface{}) interface{} {
	return &File{}
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

			id := oid.NewObjectId(0, AppConfig.NodeId()).String()
			log.Printf("Id generated: [%s][%d][%d]\n", id, AppConfig.NodeId(), oid.ObjectIdHex(id).NodeId())

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
