package config

import (
	"fmt"
	"os"
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
	ProjectName       = "REPLACEME"
)

func InitLogger() {
	homeDir := viper.Get(EnvVarProjectDir)
	logDir := fmt.Sprintf("%s/%s/logs", homeDir, ProjectName)
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		log.Fatal().Err(err).Msg("Error failed to make logDir")
	}
	logFileName := fmt.Sprintf("%s/%d.log", logDir, time.Now().Unix())
	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal().Err(err).Msg("Error opening file")
	}

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr}
	multi := zerolog.MultiLevelWriter(consoleWriter, logFile)
	log.Logger = zerolog.New(multi).Level(zerolog.TraceLevel).With().Timestamp().Caller().Logger()
	log.Info().Msg(fmt.Sprintf("Logging to %s", logFileName))
}

func LoadConfig() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting home dir")
	}
	viper.SetDefault(EnvVarLogLevel, "debug")
	viper.SetDefault(EnvVarPort, "8888")
	viper.SetDefault(EnvVarProjectDir, homeDir)
	viper.SetDefault(EnvVarReleaseMode, "false")

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix(ProjectName + "_")

	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		log.Error().Err(err).Msg("No config file loaded")
	} else {
		log.Info().Msg(fmt.Sprintf("Using config file: %s", viper.ConfigFileUsed()))
	}
	viper.AutomaticEnv()
}
