package main

import (
	config2 "github.com/rgraterol/beers-api/cmd/api/initializers/config"
	logger2 "github.com/rgraterol/beers-api/cmd/api/initializers/logger"
	server2 "github.com/rgraterol/beers-api/cmd/api/initializers/server"
)

func main() {
	initializers()
}

func initializers () {
	config2.ConfigInitializer()
	logger2.LoggerInitializer()
	server2.ServerInitializer()
}
