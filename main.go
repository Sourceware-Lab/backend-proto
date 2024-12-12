package main

import (
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/Sourceware-Lab/backend-proto/api/greeting"
	"github.com/Sourceware-Lab/backend-proto/config"
	ginLogger "github.com/gin-contrib/logger"
)

const apiVersion = "0.0.1"

type Options struct {
	Port int `help:"Port to listen on" short:"p"`
}

func (o *Options) loadFromViper() {
	o.Port = viper.GetInt(config.EnvVarPort)
}

// This is to make testing easier. We can pass a testing API interface.
func addRoutes(api huma.API) {
	log.Info().Msg("Starting loading routes")
	// Register GET /greeting/{name}
	huma.Register(api, huma.Operation{
		OperationID: "get-greeting",
		Method:      http.MethodGet,
		Path:        "/greeting/{name}",
		Summary:     "Get a greeting",
		Description: "Get a greeting for a person by name.",
		Tags:        []string{"Greetings"},
	},
		greeting.Get,
	)

	// Register POST /reviews
	huma.Register(api, huma.Operation{
		OperationID:   "post-greeting",
		Method:        http.MethodPost,
		Path:          "/greeting",
		Summary:       "Post a greeting",
		Tags:          []string{"Greetings"},
		DefaultStatus: http.StatusCreated,
	},
		greeting.Post,
	)
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

		addRoutes(api)

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
	cli := getCli()
	cli.Run()
}
