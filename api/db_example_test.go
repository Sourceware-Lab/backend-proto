package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"

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
func TestPostOrm(t *testing.T) {
	dbName := setup()
	defer teardown(dbName)
	_, api := humatest.New(t)

	AddRoutes(api)

	birthday := "2007-09-18"
	memberNumber := strconv.Itoa(rand.Intn(1000000))

	postBody := dbexample.PostBodyInputDbExample{}.Body
	postBody.Age = 25
	postBody.Name = "Jo"
	postBody.Email = "jo@example.com"
	postBody.Birthday = &birthday
	postBody.MemberNumber = &memberNumber

	resp := api.Post("/db_example/orm", postBody)
	returnBody := dbexample.PostOutputDbExample{}.Body
	err := json.Unmarshal(resp.Body.Bytes(), &returnBody)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %s", err.Error())
	}

	correctBody := dbexample.PostOutputDbExample{}.Body
	correctBody.ID = "1"

	if !cmp.Equal(returnBody, correctBody) {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
	getResp := api.Get("/db_example/orm/1")
	getRespBody := dbexample.GetOutputDbExample{}
	err = json.Unmarshal(getResp.Body.Bytes(), &getRespBody)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %s", err.Error())
	}
	if !cmp.Equal(getRespBody.Body, postBody) {
		t.Fatalf("Unexpected response: %s", getResp.Body.String())
	}
	// TODO check for the return ID and check the DB to see if it is in there OR use a get route to check for it
}
