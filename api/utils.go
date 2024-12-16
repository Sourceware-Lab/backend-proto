package api

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	dbexample "github.com/Sourceware-Lab/go-huma-gin-postgres-template/api/db_examples"
	"github.com/Sourceware-Lab/go-huma-gin-postgres-template/api/greeting"
	"github.com/Sourceware-Lab/go-huma-gin-postgres-template/api/healthcheck"
)

// AddRoutes This is to make testing easier. We can pass a testing API interface.
//
//nolint:funlen
func AddRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "healthcheck",
		Method:      http.MethodGet,
		Path:        "/healthcheck",
		Summary:     "healthcheck",
		Description: "healthcheck returns a 200 if the server is running.",
		Tags:        []string{"Healthcheck"},
	},
		healthcheck.Get,
	)

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

	huma.Register(api, huma.Operation{
		OperationID: "get-dbexample_orm",
		Method:      http.MethodGet,
		Path:        "/db_example/orm/{id}",
		Summary:     "Get to db with orm",
		Tags:        []string{"db_example"},
	},
		dbexample.GetOrm,
	)
	huma.Register(api, huma.Operation{
		OperationID: "get-dbexample_raw_sql",
		Method:      http.MethodGet,
		Path:        "/db_example/raw_sql/{id}",
		Summary:     "Get to db with raw_sql",
		Tags:        []string{"db_example"},
	},
		dbexample.GetRawSQL,
	)

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
	huma.Register(api, huma.Operation{
		OperationID:   "post-dbexample_raw_sql",
		Method:        http.MethodPost,
		Path:          "/db_example/raw_sql",
		Summary:       "Post to db with raw sql",
		Tags:          []string{"db_example"},
		DefaultStatus: http.StatusCreated,
	},
		dbexample.PostRawSQL,
	)
}
