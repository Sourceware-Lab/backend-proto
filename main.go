package main

import (
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	"github.com/danielgtaylor/huma/v2/humacli"
	ginLogger "github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	beApi "github.com/Sourceware-Lab/backend-proto/api"
	"github.com/Sourceware-Lab/backend-proto/config"
	DBpostgres "github.com/Sourceware-Lab/backend-proto/database/postgres"
)

const apiVersion = "0.0.1"

type Options struct {
	Port int `help:"Port to listen on" short:"p"`
}

func (o *Options) loadFromViper() {
	o.Port = viper.GetInt(config.EnvVarPort)
}

func getCli() (cli humacli.CLI) { // this -> (cli humacli.CLI) is a really cool go feature. It inits the var, and
	// when you use a raw return it will return the var called cli. This improves the go auto docs.

	cli = humacli.New(func(hooks humacli.Hooks, options *Options) {
		log.Info().Msg("Starting server")
		options.loadFromViper()

		if viper.Get(config.EnvVarReleaseMode) == "true" {
			gin.SetMode(gin.ReleaseMode)
		} else {
			gin.SetMode(gin.DebugMode)
		}
		gin.DisableConsoleColor()
		gin.DefaultWriter = log.Logger
		gin.DefaultErrorWriter = log.Logger

		// Create a new router & API
		router := gin.New()
		router.Use(ginLogger.SetLogger())
		api := humagin.New(router, huma.DefaultConfig("Example API", apiVersion))

		beApi.AddRoutes(api)

		// Tell the CLI how to start your server.
		hooks.OnStart(func() {
			log.Info().Msg(fmt.Sprintf("Starting server on port %d...\n", options.Port))
			_ = http.ListenAndServe(fmt.Sprintf(":%d", options.Port), router)
		})
	})
	return
}

func main() {
	config.LoadConfig()
	config.InitLogger()
	DBpostgres.Open(viper.GetString(config.EnvVarDatabaseDSN))
	DBpostgres.RunMigrations()
	cli := getCli()
	cli.Run()
}
