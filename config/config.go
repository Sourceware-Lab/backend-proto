package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
)

const (
	EnvVarLogLevel                 = "LOG_LEVEL"
	EnvVarPort                     = "PORT"
	EnvVarProjectDir               = "PROJECT_DIR"
	EnvVarReleaseMode              = "RELEASE_MODE"
	EnvVarDatabaseDSN              = "DATABASE_DSN"
	EnvVarOTELExporterOTLPEndpoint = "REPLACEME_OTEL_EXPORTER_OTLP_ENDPOINT"
)

const ProjectName = "REPLACEME"

var (
	Config config //nolint:gochecknoglobals
	Tracer = otel.Tracer("REPLACEME")
)

type config struct {
	LogLevel                 string `mapstructure:"LOG_LEVEL"`
	Port                     int    `mapstructure:"PORT"`
	ProjectDir               string `mapstructure:"PROJECT_DIR"`
	ReleaseMode              bool   `mapstructure:"RELEASE_MODE"`
	DatabaseDSN              string `mapstructure:"DATABASE_DSN"`
	OTELExporterOTLPEndpoint string `mapstructure:"REPLACEME_OTEL_EXPORTER_OTLP_ENDPOINT"`
}
type DBDSN struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
}

func (d *DBDSN) ParseDSN(dsn string) DBDSN {
	parts := make(map[string]string)

	for _, part := range strings.Split(dsn, " ") {
		kv := strings.SplitN(part, "=", 2) //nolint:mnd
		if len(kv) == 2 {                  //nolint:mnd
			parts[kv[0]] = kv[1]
		}
	}

	d.Host = parts["host"]
	d.Port, _ = strconv.Atoi(parts["port"])
	d.User = parts["user"]
	d.Password = parts["password"]
	d.DBName = parts["dbname"]
	d.SSLMode = parts["sslmode"]
	d.TimeZone = parts["TimeZone"]

	return *d
}

func (d *DBDSN) String() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		d.Host, d.User, d.Password, d.DBName, d.Port, d.SSLMode, d.TimeZone,
	)
}

func InitLogger() {
	homeDir := Config.ProjectDir
	logDir := fmt.Sprintf("%s/%s/logs", homeDir, ProjectName)

	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		log.Fatal().Err(err).Msg("Error failed to make logDir")
	}

	logFileName := fmt.Sprintf("%s/%d.log", logDir, time.Now().Unix())

	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666) //nolint:mnd
	if err != nil {
		log.Fatal().Err(err).Msg("Error opening file")
	}

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr}
	multi := zerolog.MultiLevelWriter(consoleWriter, logFile)
	log.Logger = zerolog.New(multi).Level(zerolog.TraceLevel).With().Timestamp().Caller().Logger()
	log.Info().Msg("Logging to " + logFileName)
}

func LoadConfig() {
	err := os.Setenv("TZ", "GMT")
	if err != nil {
		log.Fatal().Err(err).Msg("Error setting timezone")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting home dir")
	}

	viper.SetDefault(EnvVarLogLevel, "debug")
	viper.SetDefault(EnvVarPort, "8080")
	viper.SetDefault(EnvVarProjectDir, homeDir)
	viper.SetDefault(EnvVarReleaseMode, "false")
	viper.SetDefault(EnvVarDatabaseDSN,
		"host=localhost user=postgres password=local_fake dbname=postgres port=5432 sslmode=disable TimeZone=GMT",
	)
	viper.SetDefault(EnvVarOTELExporterOTLPEndpoint, "0.0.0.0:4317")

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	err = viper.ReadInConfig() // Find and read the config file

	if err != nil { // Handle errors reading the config file
		log.Error().Err(err).Msg("No config file loaded")
	} else {
		log.Info().Msg("Using config file: " + viper.ConfigFileUsed())
	}

	viper.AutomaticEnv()

	err = viper.Unmarshal(&Config)
	if err != nil {
		log.Fatal().Err(err).Msg("Error unmarshalling config")
	}
}
