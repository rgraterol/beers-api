package server

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"net/http/pprof"

	"github.com/rgraterol/beers-api/initializers/config"
	"github.com/rgraterol/beers-api/initializers/logger"
)

var handler http.Handler

var Server *http.ServeMux

var ServerConfig ServerConfiguration

// ServerConfiguration represents a server configuration.
type ServerConfiguration struct {
	// Address is where the Server will listen
	Address string `yaml:"address"`
	// Debug flag if we should enable debugging features.
	Debug bool `yaml:"debug"`
}

func basePingHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("pong")
}

func serverInitializer() {
	err := config.LoadConfigSection("server", &ServerConfig)
	if err != nil {
		panic(errors.New("failed to read the server config"))
	}

	Server = http.NewServeMux()
	Server.HandleFunc("/ping", basePingHandler)

	if ServerConfig.Debug {
		Server.HandleFunc("/debug/pprof/", pprof.Index)
		Server.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		Server.HandleFunc("/debug/pprof/profile", pprof.Profile)
		Server.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		Server.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}
	handler = Server
}

func RunServer() error {
	config.ConfigInitializer()
	logger.LoggerInitializer()
	serverInitializer()

	zap.S().Info("Application running on address ", ServerConfig.Address, " and enviroment ", config.Env())
	return http.ListenAndServe(ServerConfig.Address, handler)
}
