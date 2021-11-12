package router

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	log "github.com/sirupsen/logrus"

	"github.com/website/handlers/contact"
	"github.com/website/handlers/download"
	recaptcha "github.com/website/handlers/reCaptcha"
)

const (
	apiVersion1 = "/api/v1"
)

// FileSystem is a custom file system handler to handle requests to React routes
type FileSystem struct {
	fs http.FileSystem
}

// Open opens file
func (fs FileSystem) Open(path string) (http.File, error) {
	index := "/index.html"

	f, err := fs.fs.Open(path)
	if os.IsNotExist(err) {
		if f, err = fs.fs.Open(index); err != nil {
			log.Error(err)
			return nil, err
		}
	} else if err != nil {
		log.Error(err)
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if s.IsDir() {
		if _, err = fs.fs.Open(index); err != nil {
			log.Error(err)
			return nil, err
		}
	}

	return f, nil
}

//Router is the wrapper for go chi
type Router struct {
	*chi.Mux
}

//NewRouter creates new router
func NewRouter() *Router {
	r := chi.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(c.Handler)
	return &Router{Mux: r}
}

//AddRoutes adds routes to the router
func (router *Router) AddRoutes() {
	contactHandler := contact.New()
	downloadHanlder := download.New()
	reCaptchaHandler := recaptcha.New()

	router.Group(func(r chi.Router) {
		//routes to contact handler
		r.Post(apiVersion1+"/contact", contactHandler.PostContact)

		//routes to download handler
		r.Get(apiVersion1+"/download/resume", downloadHanlder.GetResume)

		//routes to reCaptcha handler
		r.Get(apiVersion1+"/recaptcha", reCaptchaHandler.GetToken)

		// paths that don't exist in the API server
		r.HandleFunc("/api/*", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			w.Write([]byte("Resource not available"))
		})

		r.HandleFunc("/googlefb60a2c818affcda.html", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "static/googlefb60a2c818affcda.html")
		})
	})
	// set up static file serving
	fs := http.FileServer(FileSystem{
		fs: http.Dir("build"),
	})
	router.Handle("/*", fs)
}
