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
	"os"
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

func (this *Geometry) FindAll(enc encoder.Encoder, params martini.Params, opt models.URLOptionsScheme, w http.ResponseWriter, r *http.Request) (int, []byte) {

	if opt.Limit == 0 || opt.Limit > 100 {
		opt.Limit = 25
	}

	result := this.model.FindAll(opt)

	return http.StatusOK, encoder.Must(enc.Encode(result))
}

func (this *Geometry) Find(enc encoder.Encoder, params martini.Params, options models.URLOptionsScheme, w http.ResponseWriter, r *http.Request) (int, []byte) {

	result := this.model.Find(params["id"])

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

func (this *Geometry) Create(payload models.GeometryScheme, enc encoder.Encoder) (int, []byte) {
	log.Println("payload:", payload)

	result, err := this.model.Create(payload)
	if err != nil {
		return http.StatusBadRequest, []byte{}
	}

	return http.StatusOK, encoder.Must(enc.Encode(result))
}

func (this *Geometry) Put(payload models.GeometryScheme, enc encoder.Encoder, p martini.Params) (int, []byte) {

	result, err := this.model.Put(p["id"], payload)
	if err != nil {
		return http.StatusBadRequest, []byte{}
	}

	return http.StatusOK, encoder.Must(enc.Encode(result))
}

func (this *Geometry) Remove(enc encoder.Encoder, p martini.Params) (int, []byte) {
	err := this.model.Remove(p["id"])
	if err != nil {
		return http.StatusBadRequest, []byte{}
	}

	return http.StatusOK, []byte{}
}
