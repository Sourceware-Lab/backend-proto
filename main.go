package main

import (
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/gin-gonic/gin"

	beApi "github.com/Sourceware-Lab/backend-proto/api"
)

const apiVersion = "0.0.1"

type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"8888"`
}

func getCli() (cli humacli.CLI) { // this -> (cli humacli.CLI) is a really cool go feature. It inits the var, and
	// when you use a raw return it will return the var called cli. This improves the go auto docs.
	cli = humacli.New(func(hooks humacli.Hooks, options *Options) {
		// Create a new router & API
		router := gin.New()
		api := humagin.New(router, huma.DefaultConfig("Example API", apiVersion))

		beApi.AddRoutes(api)

		// Tell the CLI how to start your server.
		hooks.OnStart(func() {
			fmt.Printf("Starting server on port %d...\n", options.Port)
			_ = http.ListenAndServe(fmt.Sprintf(":%d", options.Port), router)
		})
	})
	return
}

func main() {
	cli := getCli()
	cli.Run()
}
