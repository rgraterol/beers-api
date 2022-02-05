package server

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rgraterol/beers-api/initializers/config"
	"github.com/rgraterol/beers-api/initializers/logger"
)

var serverConfig ServerConfiguration

// ServerConfiguration represents a server configuration.
type ServerConfiguration struct {
	// Address is where the Server will listen
	Address string `yaml:"address"`
	// Timeout for all requests.
	Timeout int `yaml:"timeout"`
}

func basePingHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("pong")
}

func ServerInitializer() {
	err := config.LoadConfigSection("server", &serverConfig)
	if err != nil {
		panic(errors.New("failed to read the server config"))
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Duration(serverConfig.Timeout) * time.Second))
	r.Use(logger.ChiLogger())

	r.Get("/ping", basePingHandler)

	zap.S().Info("Application running on address ", serverConfig.Address, " and enviroment ", config.Env())
	http.ListenAndServe(serverConfig.Address, r)
}
