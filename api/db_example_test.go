package api

import (
	"fmt"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/uuid"

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
	payload := strings.NewReader("{\n  \"Age\": 25,\n  \"Birthday\": \"2006-01-02\",\n  \"Email\": \"jo@example.com\",\n  \"MemberNumber\": \"123456\",\n  \"Name\": \"Jo\"\n}")
	resp := api.Post("/db_example/orm", payload)
	if !strings.Contains(resp.Body.String(), "{\"id\":\"1\"}") { // TODO check for the return ID and check the DB to see if it is in there
		// OR use a get route to check for it
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}
