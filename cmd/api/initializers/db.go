package initializers

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

var DB *sql.DB

var DatabaseConfig DatabaseConfiguration

// DatabaseConfiguration represents a database configuration.
type DatabaseConfiguration struct {
	// URL is the database address.
	URL string `yaml:"url"`
	// Debug flag if we should enable debugging features.
	Debug bool `yaml:"debug"`
}

func DatabaseInitializer() {
	err := LoadConfigSection("database", &DatabaseConfig)
	if err != nil {
		panic(errors.Wrap(err, "failed to read the database config"))
	}

	if url := os.Getenv("DATABASE_URL"); url != "" {
		DatabaseConfig.URL = url
	}

	DB, err = sql.Open("mysql", DatabaseConfig.URL)
	if err != nil {
		panic(errors.Wrap(err, "failed to initialize the DB"))
	}
}
