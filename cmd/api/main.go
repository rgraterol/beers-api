package main

import (
	"github.com/rgraterol/beers-api/initializers/config"
	"github.com/rgraterol/beers-api/initializers/logger"
	"github.com/rgraterol/beers-api/initializers/server"
)

func main() {
	initializers()
}

func initializers () {
	config.ConfigInitializer()
	logger.LoggerInitializer()
	server.ServerInitializer()
}
