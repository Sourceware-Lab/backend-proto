package DBpostgres

import (
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Open(dsn string) {
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
	retries := 3
	retry := 0
	var db *gorm.DB
	var err error
	for retry < retries {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
		if err == nil {
			break
		}
		log.Info().Err(err).Msg("Error connecting to database, retrying in 3 seconds")
		retry += 1
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to database")
	}

	DB = db
}

func RunMigrations() {
	err := DB.AutoMigrate(&User{})
	if err != nil {
		log.Fatal().Err(err).Msg("Error migrating database")
	}
}
