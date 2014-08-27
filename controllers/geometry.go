package controllers

import (
	// "df/api/utils"
	"df/assets/models"
	. "df/assets/utils"
	"github.com/martini-contrib/encoder"
	// "io/ioutil"
	"github.com/go-martini/martini"
	"log"
	"net/http"
	// "os"
	"strings"
)

type Geometry struct {
	model *models.Geometry
}

func (*Geometry) Construct(args ...interface{}) interface{} {
	return &Geometry{
		model: (*models.Geometry).Construct(nil).(*models.Geometry),
	}
}

func (this *Geometry) Find(enc encoder.Encoder, params martini.Params, options models.URLOptionsScheme, w http.ResponseWriter, r *http.Request) (int, []byte) {
	var httpStatus int = http.StatusOK

	result := this.model.Find(params["id"])
	if result == nil {
		httpStatus = http.StatusNotFound
	}

	if strings.Contains("application/json", r.Header.Get("Accept")) {
		return httpStatus, encoder.Must(enc.Encode(result))
	}

	if result == nil || result.Base == "" {
		return httpStatus, []byte{}
	}

	pmap := options.ToMap()
	// name will be /path/to/id + optionally ".{param}{param}{param}{param}"
	name := AppConfig.StorageFilePath(result.Base) + options.ToHash(pmap)
	log.Println("Serving file:", name)

	if err := result.Morph(name, pmap); err != nil {
		log.Println("Unable to morph:", err)
		return http.StatusNotFound, []byte{}
	}

	w.Header().Set("Content-Type", "application/octet-stream")

	http.ServeFile(w, r, name)

	return http.StatusOK, []byte{}
}

func (this *Geometry) Create(payload models.GeometryScheme, enc encoder.Encoder) (int, []byte) {
	log.Println("payload:", payload)

	result, err := this.model.Create(payload)
	if err != nil {
		return http.StatusBadRequest, []byte{}
	}

	return http.StatusOK, encoder.Must(enc.Encode(result))
}
