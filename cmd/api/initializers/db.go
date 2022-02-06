package initializers

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"path"
	"time"

	"github.com/pkg/errors"
	"github.com/rgraterol/beers-api/pkg/db"
	"github.com/rgraterol/beers-api/pkg/usecases/beers"
)

var DatabaseConfig DatabaseConfiguration

// DatabaseConfiguration represents a database configuration.
type DatabaseConfiguration struct {
	// URL is the database address.
	URL string `yaml:"url"`
	// MaxIdleConns sets the maximum number of connections in the idle connection pool.
	MaxIdleConns int `yaml:"maxIdleConns"`
	// MaxOpenConns sets the maximum number of open connections to the database.
	MaxOpenConns int `yaml:"maxOpenConns"`
	// ConnMaxLifetime sets the maximum amount of time in minutes a connection may be reused.
	ConnMaxLifetime int `yaml:"connMaxLifetime"`
	// Automigrate set condition to automatically migrate db schema.
	AutoMigrate bool `yaml:"autoMigrate"`
}

func DatabaseInitializer() {
	err := LoadConfigSection("database", &DatabaseConfig)
	if err != nil {
		panic(errors.Wrap(err, "failed to read the database config"))
	}

	if url := os.Getenv("DATABASE_URL"); url != "" {
		DatabaseConfig.URL = url
	}

	db.Gorm, err = gorm.Open(mysql.Open(DatabaseConfig.URL), &gorm.Config{Logger: initGormLogger()})
	if err != nil {
		panic(errors.Wrap(err, "failed to initialize the DB"))
	}
	pool, err := db.Gorm.DB()
	if err != nil {
		panic(errors.Wrap(err, "failed to configure connection pool"))
	}
	pool.SetMaxIdleConns(DatabaseConfig.MaxIdleConns)
	pool.SetMaxOpenConns(DatabaseConfig.MaxOpenConns)
	pool.SetConnMaxLifetime(time.Duration(DatabaseConfig.ConnMaxLifetime))

	if DatabaseConfig.AutoMigrate {
		err = runMigrations()
		if err != nil {
			panic(err)
		}
	}
}

func MockDatabaseInitializer() {
	testDBPath := rootDir() + "/db/mock/test.db"
	err := clearTestDB(testDBPath)
	if err != nil {
		panic(errors.Wrap(err, "cannot clear test DB"))
	}

	db.Gorm, err = gorm.Open(sqlite.Open(testDBPath), &gorm.Config{Logger: initGormLogger()})
	if err != nil {
		panic(errors.Wrap(err, "failed to connect gorm with mock DB"))
	}
	err = runMigrations()
	if err != nil {
		panic(err)
	}
}

func rootDir() string {
	crrPath, _ := crrFSGetter.getwd()
	dir, err := os.Open(path.Join(crrPath, "../../../"))
	if err != nil {
		panic(errors.New("cannot find project root directory"))
	}
	return dir.Name()
}

func clearTestDB(path string) error {
	if _, err := os.Stat(path); err == nil {
		return os.Remove(path)
	}
	return nil
}

func runMigrations() error {
	if db.Gorm.Migrator().HasTable(&beers.Beer{}) {
		return nil
	}
	err := db.Gorm.AutoMigrate(&beers.Beer{})
	if err != nil {
		return errors.Wrap(err,  "cannot run beers migration")
	}
	return nil
}

func initGormLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:              time.Second,
			LogLevel:                   logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
}
