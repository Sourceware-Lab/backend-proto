//go:build integration
// +build integration

package api

import (
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
)

func setup() {
	// TODO setup the DB
}
func TestPostOrm(t *testing.T) {
	setup()
	_, api := humatest.New(t)

	AddRoutes(api)
	payload := strings.NewReader("{\n  \"Age\": 25,\n  \"Birthday\": \"2006-01-02\",\n  \"Email\": \"jo@example.com\",\n  \"MemberNumber\": \"123456\",\n  \"Name\": \"Jo\"\n}")
	resp := api.Post("/db_example/orm", payload)
	if !strings.Contains(resp.Body.String(), "Hello post, test!") { // TODO check for the return ID and check the DB to see if it is in there
		// OR use a get route to check for it
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}
