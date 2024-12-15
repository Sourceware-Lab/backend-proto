package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"

	dbexample "github.com/Sourceware-Lab/backend-proto/api/db_examples"
	"github.com/Sourceware-Lab/backend-proto/config"
	DBpostgres "github.com/Sourceware-Lab/backend-proto/database/postgres"
)

func setup() (dbName string) {
	config.LoadConfig()

	DBpostgres.Open(config.Config.DatabaseDSN)

	dbDSNString := config.Config.DatabaseDSN
	dbDSN := config.DbDSN{}
	dbDSN.ParseDSN(dbDSNString)
	dbName = strings.Replace(fmt.Sprintf("testdb-%s", uuid.New().String()), "-", "", -1)
	dbDSN.DbName = dbName

	DBpostgres.CreateDb(dbName)

	DBpostgres.Open(dbDSN.String())
	DBpostgres.RunMigrations()
	return dbName
}

func teardown(dbName string) {
	DBpostgres.Open(config.Config.DatabaseDSN)
	DBpostgres.DeleteDb(dbName)
}

//nolint:funlen
func TestRoutes(t *testing.T) {
	tests := []struct {
		name     string
		basePath string
		want     dbexample.PostInputDbExample
	}{
		{
			name:     "get",
			basePath: "/db_example/orm",
			want: dbexample.PostInputDbExample{
				Body: dbexample.PostBodyInputDbExampleBody{
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
			want: dbexample.PostInputDbExample{
				Body: dbexample.PostBodyInputDbExampleBody{
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
			_, api := humatest.New(t)
			AddRoutes(api)
			birthdayTime := time.Now().Add(time.Duration(-tt.want.Body.Age) * time.Hour * 24 * 365)
			birthday := birthdayTime.Format(time.DateOnly)
			tt.want.Body.Birthday = &birthday

			memberNumber := strconv.Itoa(rand.Intn(1000000))
			tt.want.Body.MemberNumber = &memberNumber

			resp := api.Post(tt.basePath, tt.want.Body)

			postRespBody := dbexample.PostOutputDbExample{}.Body
			err := json.Unmarshal(resp.Body.Bytes(), &postRespBody)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %s", err.Error())
			}

			expectedPostBody := dbexample.PostOutputDbExample{}.Body
			expectedPostBody.ID = "1"
			if !cmp.Equal(postRespBody, expectedPostBody) {
				t.Fatalf("Unexpected response: %s", resp.Body.String())
			}

			getResp := api.Get(tt.basePath + "/1")
			getRespBody := dbexample.GetOutputDbExample{}
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
