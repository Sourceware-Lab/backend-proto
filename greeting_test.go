package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"

	"github.com/Sourceware-Lab/backend-proto/api/greeting"
)

func TestGetGreeting(t *testing.T) {
	_, api := humatest.New(t)

	addRoutes(api)

	resp := api.Get("/greeting/world")
	if !strings.Contains(resp.Body.String(), "Hello get, world!") {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}

func TestPostGreeting(t *testing.T) {
	_, api := humatest.New(t)

	addRoutes(api)

	resp := api.Post("/greeting",
		map[string]any{
			"name": "test",
		},
	)
	if !strings.Contains(resp.Body.String(), "Hello post, test!") {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}

func TestPostMissingBody(t *testing.T) {
	_, api := humatest.New(t)

	addRoutes(api)

	resp := api.Post("/greeting",
		map[string]any{
			"FAKE": "test",
		},
	)

	if resp.Code != 422 {
		t.Fatalf("Unexpected status code: %d", resp.Code)
	}
}

func FuzzPostGreeting(f *testing.F) {
	f.Add("hello")
	f.Fuzz(
		func(t *testing.T, name string) {
			if len(name) >= 30 {
				name = name[:29] // Truncate the string to ensure it's less than 30 because that is a limit of the endpoint
			}
			_, api := humatest.New(t)
			addRoutes(api)
			resp := api.Post("/greeting",
				map[string]any{
					"name": name,
				},
			)
			var jsonData greeting.Output
			json.Unmarshal([]byte(resp.Body.String()), &jsonData.Body)
			if jsonData.Body.Message != fmt.Sprintf("Hello post, %s!", name) {
				t.Fatalf("Unexpected response: %s", resp.Body.String())
			}
		},
	)
}
