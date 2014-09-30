package main

import (
	ctrl "df/assets/controllers"
	"df/assets/models"
	. "df/assets/utils"
	"fmt"
	"github.com/3d0c/binding"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/encoder"
	"log"
	"net/http"
	"runtime"

	// "flag"
	// "os"
	// "os/signal"
	// "runtime/pprof"
	// "syscall"
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	InitConfigFrom("./config.json")
	runtime.GOMAXPROCS(8)
}

func main() {
	/*
		var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

		flag.Parse()

		if *cpuprofile != "" {
			f, err := os.Create(*cpuprofile)
			if err != nil {
				log.Fatal(err)
			}
			pprof.StartCPUProfile(f)
		}

		signalChannel := make(chan os.Signal, 2)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
		go func() {
			sig := <-signalChannel
			switch sig {
			case os.Interrupt:
				log.Println("Interrupted")
				pprof.StopCPUProfile()
				os.Exit(0)
			case syscall.SIGTERM:
				log.Println("Terminated")
				pprof.StopCPUProfile()
				os.Exit(0)
			}
		}()
	*/
	m := martini.New()
	route := martini.NewRouter()
	// pr := martini.Recovery()

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

	route.Get("/geometry",
		binding.Bind(models.URLOptionsScheme{}),
		construct(&ctrl.Geometry{}),
		(*ctrl.Geometry).FindAll,
	)

	route.Post("/geometry",
		binding.Bind(models.GeometryScheme{}),
		construct(&ctrl.Geometry{}),
		(*ctrl.Geometry).Create,
	)

	route.Put("/geometry/:id",
		binding.Bind(models.GeometryScheme{}),
		construct(&ctrl.Geometry{}),
		(*ctrl.Geometry).Put,
	)

	route.Delete(
		"/geometry/:id",
		construct(&ctrl.Geometry{}),
		(*ctrl.Geometry).Remove,
	)

	route.Get("/:id",
		construct(&ctrl.File{}),
		(*ctrl.File).Find,
	)

	route.Post("/",
		construct(&ctrl.File{}),
		(*ctrl.File).Create,
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
