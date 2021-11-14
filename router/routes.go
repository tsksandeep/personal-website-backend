package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"

	anytocsv "github.com/website/handlers/anyToCsv"
	"github.com/website/handlers/contact"
	"github.com/website/handlers/download"
	recaptcha "github.com/website/handlers/reCaptcha"
)

const (
	apiVersion1 = "/api/v1"
)

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
	anyToCsvHandler := anytocsv.New()

	router.Group(func(r chi.Router) {
		//routes to contact handler
		r.Post(apiVersion1+"/contact", contactHandler.PostContact)

		//routes to download handler
		r.Get(apiVersion1+"/download/resume", downloadHanlder.GetResume)

		//routes to reCaptcha handler
		r.Get(apiVersion1+"/recaptcha", reCaptchaHandler.GetToken)

		//routes to anyToCsv handler
		r.Post(apiVersion1+"/anytocsv", anyToCsvHandler.PostCsv)

		// paths that don't exist in the API server
		r.HandleFunc("/api/*", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			w.Write([]byte("Resource not available"))
		})
	})
}
