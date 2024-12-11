package main

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	"github.com/gin-gonic/gin"

	"github.com/Sourceware-Lab/backend-proto/api/greeting"
)

// GreetingOutput represents the greeting operation response.

func addRoutes(api huma.API) {
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

func main() {
	// Create a new router & API
	router := gin.New()
	api := humagin.New(router, huma.DefaultConfig("My API", "0.0.1"))

	addRoutes(api)

	http.ListenAndServe("0.0.0.0:4000", router)
}
