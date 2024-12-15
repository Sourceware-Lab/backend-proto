package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"

	ginLogger "github.com/gin-contrib/logger"

	beApi "github.com/Sourceware-Lab/backend-proto/api"
	"github.com/Sourceware-Lab/backend-proto/config"
	DBpostgres "github.com/Sourceware-Lab/backend-proto/database/postgres"
)

const apiVersion = "0.0.1"

type Options struct {
	Port int `help:"Port to listen on" short:"p"`
}

func (o *Options) loadFromViper() {
	o.Port = config.Config.Port
}

//nolint:ireturn
func getCli() (cli humacli.CLI) { // this -> (cli humacli.CLI) is a really cool go feature. It inits the var, and
	// when you use a raw return it will return the var called cli. This improves the go auto docs.

	cli = humacli.New(func(hooks humacli.Hooks, options *Options) {
		log.Info().Msg("Starting server")
		options.loadFromViper()

		if config.Config.ReleaseMode {
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
			server := &http.Server{
				IdleTimeout:       300 * time.Second, //nolint:mnd
				ReadTimeout:       300 * time.Second, //nolint:mnd
				WriteTimeout:      300 * time.Second, //nolint:mnd
				ReadHeaderTimeout: 10 * time.Second,  //nolint:mnd
				Addr:              fmt.Sprintf(":%d", options.Port),
				Handler:           router,
			}
			_ = server.ListenAndServe()
		})
	})

	return cli
}

func main() {
	config.LoadConfig()
	config.InitLogger()
	DBpostgres.Open(config.Config.DatabaseDSN)
	defer DBpostgres.Close()
	DBpostgres.RunMigrations()
	cli := getCli()
	cli.Run()
}
