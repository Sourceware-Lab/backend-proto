package dbpostgres

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/Sourceware-Lab/backend-proto/config"
)

var DB *gorm.DB //nolint:gochecknoglobals

func Open(dsn string) {
	if DB != nil {
		openDb, err := DB.DB()
		if err != nil {
			log.Fatal().Err(err).Msg("Error getting DB")
		}
		err = openDb.Close()
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
	for retry < retries {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
		if err == nil {
			break
		}
		log.Info().Err(err).Msg("Error connecting to database, retrying in 3 seconds")
		retry++
		time.Sleep(3 * time.Second)
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
}

func Close() {
	openDb, err := DB.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting DB")
	}
	err = openDb.Close()
	if err != nil {
		log.Fatal().Err(err).Msg("Error closing DB")
	}
	DB = nil
}

func CreateDb(dbName string) {
	result := DB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if result.Error != nil {
		log.Fatal().Err(result.Error).Msg("Error creating database")
	}
}

func DeleteDb(dbName string) {
	result := DB.Exec(fmt.Sprintf("DROP DATABASE %s", dbName))
	if result.Error != nil {
		log.Fatal().Err(result.Error).Msg("Error deleting database")
	}
}

func RunMigrations() {
	err := DB.AutoMigrate(&User{})
	if err != nil {
		log.Fatal().Err(err).Msg("Error migrating database")
	}
}
