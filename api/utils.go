package api

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	dbexample "github.com/Sourceware-Lab/backend-proto/api/db_examples"
	"github.com/Sourceware-Lab/backend-proto/api/greeting"
)

// AddRoutes This is to make testing easier. We can pass a testing API interface.
func AddRoutes(api huma.API) {
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

	//huma.Register(api, huma.Operation{
	//	OperationID:   "post-dbexample_raw",
	//	Method:        http.MethodPost,
	//	Path:          "/db_example/raw",
	//	Summary:       "Post to db with raw sql",
	//	Tags:          []string{"db_example"},
	//	DefaultStatus: http.StatusCreated,
	//},
	//	dbexample.PostRawSql,
	//)
	huma.Register(api, huma.Operation{
		OperationID:   "post-dbexample_orm",
		Method:        http.MethodPost,
		Path:          "/db_example/orm",
		Summary:       "Post to db with orm",
		Tags:          []string{"db_example"},
		DefaultStatus: http.StatusCreated,
	},
		dbexample.PostOrm,
	)
}
