package dbexample_test

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"

	"github.com/Sourceware-Lab/go-huma-gin-postgres-template/api"
	dbexample "github.com/Sourceware-Lab/go-huma-gin-postgres-template/api/db_examples"
	"github.com/Sourceware-Lab/go-huma-gin-postgres-template/config"
	DBpostgres "github.com/Sourceware-Lab/go-huma-gin-postgres-template/database/postgres"
)

func setup() string {
	config.LoadConfig()

	DBpostgres.Open(config.Config.DatabaseDSN)

	dbDSNString := config.Config.DatabaseDSN
	dbDSN := config.DBDSN{}
	dbDSN.ParseDSN(dbDSNString)

	dbName := strings.ReplaceAll("testdb-"+uuid.New().String(), "-", "")
	dbDSN.DBName = dbName

	DBpostgres.CreateDB(dbName)

	DBpostgres.Open(dbDSN.String())
	DBpostgres.RunMigrations()

	return dbName
}

func teardown(dbName string) {
	DBpostgres.Open(config.Config.DatabaseDSN)
	DBpostgres.DeleteDB(dbName)
}

//nolint:funlen
func TestRoutes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		basePath string
		want     dbexample.PostInputDBExample
	}{
		{
			name:     "get",
			basePath: "/db_example/orm",
			want: dbexample.PostInputDBExample{
				Body: dbexample.PostBodyInputDBExampleBody{
					Name:         "jo",
					Age:          25,
					Email:        "jo@example.com",
					Birthday:     nil,
					MemberNumber: nil,
				},
			},
		},
		{
			name:     "get",
			basePath: "/db_example/raw_sql",
			want: dbexample.PostInputDBExample{
				Body: dbexample.PostBodyInputDBExampleBody{
					Name:         "jo1",
					Age:          26,
					Email:        "jo1@example.com",
					Birthday:     nil,
					MemberNumber: nil,
				},
			},
		},
	}

	for _, tt := range tests {
		func() {
			dbName := setup()
			defer teardown(dbName)

			_, apiInstance := humatest.New(t)
			api.AddRoutes(apiInstance)

			birthdayTime := time.Now().Add(time.Duration(-tt.want.Body.Age) * time.Hour * 24 * 365)
			birthday := birthdayTime.Format(time.DateOnly)
			tt.want.Body.Birthday = &birthday

			memberNumber := strconv.Itoa(rand.Intn(1000000)) //nolint:gosec
			tt.want.Body.MemberNumber = &memberNumber

			resp := apiInstance.Post(tt.basePath, tt.want.Body)

			postRespBody := dbexample.PostOutputDBExample{}.Body

			err := json.Unmarshal(resp.Body.Bytes(), &postRespBody)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %s", err.Error())
			}

			expectedPostBody := dbexample.PostOutputDBExample{}.Body
			expectedPostBody.ID = "1"

			if !cmp.Equal(postRespBody, expectedPostBody) {
				t.Fatalf("Unexpected response: %s", resp.Body.String())
			}

			getResp := apiInstance.Get(tt.basePath + "/1")
			getRespBody := dbexample.GetOutputDBExample{}

			err = json.Unmarshal(getResp.Body.Bytes(), &getRespBody)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %s", err.Error())
			}

			if !cmp.Equal(getRespBody.Body, tt.want.Body) {
				t.Fatalf("Unexpected response: %s", getResp.Body.String())
			}
		}()
	}
}
