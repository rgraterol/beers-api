package main

import (
	i "github.com/rgraterol/beers-api/cmd/api/initializers"
)

func main() {
	run()
}

func run() {
	i.ConfigInitializer()
	i.LoggerInitializer()
	i.DatabaseInitializer()
	i.RestClientsInitializer()
	i.ServerInitializer()
}
