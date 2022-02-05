package initializers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rgraterol/beers-api/pkg/beers"
)

func routes(r *chi.Mux) {
	r.Get("/ping", basePingHandler)

	r.Route("/beers", func(r chi.Router) {
		r.Get("/", beers.List)
	})
}

func basePingHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("pong")
}
