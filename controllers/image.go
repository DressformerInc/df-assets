package controllers

import (
	// "df/api/utils"
	"df/assets/models"
	. "df/assets/utils"
	// "github.com/3d0c/oid"
	img "github.com/3d0c/imgproc"
	"github.com/martini-contrib/encoder"
	// "io/ioutil"
	"github.com/go-martini/martini"
	"log"
	"net/http"
	// "os"
	// "strings"
)

type Image struct{}

func (*Image) Construct(args ...interface{}) interface{} {
	return &Image{}
}

func (this *Image) Find(enc encoder.Encoder, params martini.Params, options models.URLOptionsScheme, w http.ResponseWriter) (int, []byte) {
	log.Println(params, options)
	options = options.SetDefaults()

	base := img.NewSource(AppConfig.StorageFilePath(params["id"]))
	if base == nil {
		return http.StatusNotFound, []byte{}
	}

	target := &img.Options{
		Base:    base,
		Scale:   img.NewScale(options.Scale),
		Format:  options.Format,
		Method:  3,
		Quality: options.Quality,
	}

	w.Header().Set("Content-Type", "image/"+target.Format)
	log.Println(target)

	return http.StatusOK, img.Proc(target)
}
