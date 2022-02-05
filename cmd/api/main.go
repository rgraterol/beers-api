package main

import (
	"log"

	"github.com/rgraterol/beers-api/initializers/server"
)


func main() {
	log.Panic(server.RunServer())
}
