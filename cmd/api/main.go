package main

import (
	"github.com/rgraterol/beers-api/cmd/api/initializers"
)

func main() {
	run()
}

func run() {
	initializers.ConfigInitializer()
	initializers.LoggerInitializer()
	initializers.ServerInitializer()
}
