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
		var b beers.Service
		r.Get("/", beers.List(&b))
		r.Post("/", beers.Create(&b))
		r.Get("/{beerID}", beers.Get(&b))
	})
}

func basePingHandler(w http.ResponseWriter, _ *http.Request) {
	responses.OK(w, "pong")
}
