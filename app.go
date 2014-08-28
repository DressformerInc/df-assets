package main

import (
	ctrl "df/assets/controllers"
	"df/assets/models"
	. "df/assets/utils"
	"fmt"
	"github.com/3d0c/binding"
	"github.com/3d0c/martini-contrib/config"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/encoder"
	"log"
	"net/http"
	"runtime"
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	config.Init("./config.json")
	config.LoadInto(AppConfig)

	runtime.GOMAXPROCS(8)
}

func main() {
	m := martini.New()
	route := martini.NewRouter()

	m.Use(func(c martini.Context, w http.ResponseWriter) {
		c.MapTo(encoder.JsonEncoder{PrettyPrint: true}, (*encoder.Encoder)(nil))
		w.Header().Set("Content-Type", "application/json")
	})

	m.Use(func(w http.ResponseWriter, req *http.Request) {
		if origin := req.Header.Get("Origin"); origin != "" {
			w.Header().Add("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Add("Access-Control-Allow-Origin", "*")
		}

		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, X-Requested-With")
		w.Header().Add("Cache-Control", "max-age=2592000")
		w.Header().Add("Pragma", "public")
		w.Header().Add("Cache-Control", "public")
	})

	route.Options("/**")

	route.Post("/",
		construct(&ctrl.File{}),
		(*ctrl.File).Create,
	)

	route.Get("/image/:id",
		binding.Bind(models.URLOptionsScheme{}),
		construct(&ctrl.Image{}),
		(*ctrl.Image).Find,
	)

	route.Get("/geometry/:id",
		binding.Bind(models.URLOptionsScheme{}),
		construct(&ctrl.Geometry{}),
		(*ctrl.Geometry).Find,
	)

	route.Post("/geometry",
		binding.Bind(models.GeometryScheme{}),
		construct(&ctrl.Geometry{}),
		(*ctrl.Geometry).Create,
	)

	m.Action(route.Handle)

	log.Printf("Waiting for connections on %v...\n", AppConfig.ListenOn())

	go func() {
		if err := http.ListenAndServe(AppConfig.ListenOn(), m); err != nil {
			log.Fatal(err)
		}
	}()

	if err := http.ListenAndServeTLS(AppConfig.HttpsOn(), AppConfig.SSLCert(), AppConfig.SSLKey(), m); err != nil {
		log.Fatal(err)
	}
}

func construct(obj interface{}, args ...interface{}) martini.Handler {
	return func(ctx martini.Context, r *http.Request) {
		switch t := obj.(type) {
		case models.Model:
			ctx.Map(obj.(models.Model).Construct(args))

		case ctrl.Controller:
			ctx.Map(obj.(ctrl.Controller).Construct(args))

		default:
			panic(fmt.Sprintln("Unexpected type:", t))
		}
	}
}
