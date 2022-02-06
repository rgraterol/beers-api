package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rgraterol/beers-api/pkg/responses"
	"github.com/rgraterol/beers-api/pkg/usecases/beers"
)

func Routes(r *chi.Mux) {
	r.Get("/ping", basePingHandler)

	r.Route("/beers", func(r chi.Router) {
		r.Get("/", beers.List(&beers.Service{}))
		r.Post("/", beers.Create(&beers.Service{}))
	})
}

func basePingHandler(w http.ResponseWriter, _ *http.Request) {
	responses.OK(w, "pong")
}
