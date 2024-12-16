package dbpostgres

import (
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/Sourceware-Lab/go-huma-gin-postgres-template/config"
)

var DB *gorm.DB //nolint:gochecknoglobals

//nolint:funlen
func Open(dsn string) {
	log.Info().Msg("Opening database")

	if DB != nil {
		openDB, err := DB.DB()
		if err != nil {
			log.Fatal().Err(err).Msg("Error getting DB")
		}

		err = openDB.Close()
		if err != nil {
			log.Fatal().Err(err).Msg("Error closing DB")
		}

		DB = nil
	}

	dbZlog := log.Logger
	newLogger := logger.New(
		&dbZlog, // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)

	var db *gorm.DB

	var err error

	retries := 3

	retry := 0

	log.Info().Msg("Connecting to database")

	for retry < retries {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
		if err == nil {
			break
		}

		log.Info().Err(err).Msg("Error connecting to database, retrying in 3 seconds")

		retry++

		time.Sleep(3 * time.Second) //nolint:mnd
	}

	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to database")
	}

	if db == nil {
		log.Fatal().Msg("Error connecting to database")
	}

	if config.Config.ReleaseMode {
		DB = db
	} else {
		DB = db.Debug()
	}

	log.Info().Msg("Connected to database")
}

func Close() {
	openDB, err := DB.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting DB")
	}

	err = openDB.Close()
	if err != nil {
		log.Fatal().Err(err).Msg("Error closing DB")
	}

	DB = nil
}

func CreateDB(dbName string) {
	result := DB.Exec("CREATE DATABASE " + dbName)
	if result.Error != nil {
		log.Fatal().Err(result.Error).Msg("Error creating database")
	}
}

func DeleteDB(dbName string) {
	result := DB.Exec("DROP DATABASE " + dbName)
	if result.Error != nil {
		log.Fatal().Err(result.Error).Msg("Error deleting database")
	}
}

func RunMigrations() {
	log.Info().Msg("Running migrations")

	err := DB.AutoMigrate(&User{})
	if err != nil {
		log.Fatal().Err(err).Msg("Error migrating database")
	}

	log.Info().Msg("Migrations ran successfully")
}
