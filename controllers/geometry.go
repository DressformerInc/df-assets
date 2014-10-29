package controllers

import (
	"df/assets/models"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/encoder"
	"log"
	"net/http"
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

	result := this.model.FindAll([]string{}, opt)

	return http.StatusOK, encoder.Must(enc.Encode(result))
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
