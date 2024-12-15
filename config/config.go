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
)

const (
	EnvVarLogLevel    = "LOG_LEVEL"
	EnvVarPort        = "PORT"
	EnvVarProjectDir  = "PROJECT_DIR"
	EnvVarReleaseMode = "RELEASE_MODE"
	EnvVarDatabaseDSN = "DATABASE_DSN"
)

const ProjectName = "REPLACEME"

var Config config

type config struct {
	LogLevel    string `mapstructure:"LOG_LEVEL"`
	Port        int    `mapstructure:"PORT"`
	ProjectDir  string `mapstructure:"PROJECT_DIR"`
	ReleaseMode bool   `mapstructure:"RELEASE_MODE"`
	DatabaseDSN string `mapstructure:"DATABASE_DSN"`
}
type DbDSN struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
	SSLMode  string
	TimeZone string
}

func (d *DbDSN) ParseDSN(dsn string) DbDSN {
	parts := make(map[string]string)
	for _, part := range strings.Split(dsn, " ") {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			parts[kv[0]] = kv[1]
		}
	}

	d.Host = parts["host"]
	d.Port, _ = strconv.Atoi(parts["port"])
	d.User = parts["user"]
	d.Password = parts["password"]
	d.DbName = parts["dbname"]
	d.SSLMode = parts["sslmode"]
	d.TimeZone = parts["TimeZone"]

	return *d
}

func (d *DbDSN) String() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		d.Host, d.User, d.Password, d.DbName, d.Port, d.SSLMode, d.TimeZone,
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
	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatal().Err(err).Msg("Error opening file")
	}

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr}
	multi := zerolog.MultiLevelWriter(consoleWriter, logFile)
	log.Logger = zerolog.New(multi).Level(zerolog.TraceLevel).With().Timestamp().Caller().Logger()
	log.Info().Msg(fmt.Sprintf("Logging to %s", logFileName))
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
	viper.SetDefault(EnvVarPort, "8888")
	viper.SetDefault(EnvVarProjectDir, homeDir)
	viper.SetDefault(EnvVarReleaseMode, "false")
	viper.SetDefault(EnvVarDatabaseDSN, "host=localhost user=postgres password=local_fake dbname=postgres port=5432 sslmode=disable TimeZone=GMT")

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix(ProjectName)

	err = viper.ReadInConfig() // Find and read the config file

	if err != nil { // Handle errors reading the config file
		log.Error().Err(err).Msg("No config file loaded")
	} else {
		log.Info().Msg(fmt.Sprintf("Using config file: %s", viper.ConfigFileUsed()))
	}
	viper.AutomaticEnv()
	err = viper.Unmarshal(&Config)
	if err != nil {
		log.Fatal().Err(err).Msg("Error unmarshalling config")
	}
}
