package greeting_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"

	"github.com/Sourceware-Lab/go-huma-gin-postgres-template/api"
	"github.com/Sourceware-Lab/go-huma-gin-postgres-template/api/greeting"
)

// Happy path test.
func TestGetGreeting(t *testing.T) {
	t.Parallel()
	_, apiInstance := humatest.New(t)

	api.AddRoutes(apiInstance)

	resp := apiInstance.Get("/greeting/world")
	if !strings.Contains(resp.Body.String(), "Hello get, world!") {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}

func TestGetGreetingMissingPath(t *testing.T) {
	t.Parallel()
	_, apiInstance := humatest.New(t)

	api.AddRoutes(apiInstance)

	resp := apiInstance.Get("/greeting/")
	if resp.Code != 404 {
		t.Fatalf("Page found when it should not have: %s", resp.Body.String())
	}
}

// Happy path test.
func TestPostGreeting(t *testing.T) {
	t.Parallel()
	_, apiInstance := humatest.New(t)

	api.AddRoutes(apiInstance)

	resp := apiInstance.Post("/greeting",
		map[string]any{
			"name": "test",
		},
	)
	if !strings.Contains(resp.Body.String(), "Hello post, test!") {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}

func TestPostMissingBody(t *testing.T) {
	t.Parallel()

	_, apiInstance := humatest.New(t)

	api.AddRoutes(apiInstance)

	resp := apiInstance.Post("/greeting",
		map[string]any{
			"FAKE": "test",
		},
	)

	if resp.Code != 422 {
		t.Fatalf("Unexpected status code: %d", resp.Code)
	}
}

func FuzzPostGreeting(f *testing.F) {
	// These are inputs that will always be run. The first seeds. After these 3 the fuzzing lib will generate more.
	f.Add("hello")
	f.Add("\ufffd")
	f.Add("\"")

	f.Fuzz(
		func(t *testing.T, name string) {
			// Truncate the string to ensure it's less than 30 because that is a limit of the endpoint
			if len(name) >= 30 {
				name = name[:29]
			}

			_, apiInstance := humatest.New(t)
			api.AddRoutes(apiInstance)

			// Go does some fun stuff with unicode in json marshal/unmarshal. It is easier to marshal and unmarshal the
			// input data than handle the unicode correctly.
			marshaledFuzzyInput, err := json.Marshal(
				map[string]any{
					"name": name,
				},
			)
			if err != nil {
				t.Fatalf("Failed to marshal json: %s", err.Error())
			}

			// api.Post() only accepts a io.Reader interface or map, not the raw string
			resp := apiInstance.Post("/greeting", strings.NewReader(string(marshaledFuzzyInput)))

			// json unmarshal needs a map/struct to put the data, it does not return a new object.
			var unmarshalledFuzzingInput greeting.PostInputGreeting

			var jsonData greeting.OutputGreeting

			err = json.Unmarshal(marshaledFuzzyInput, &unmarshalledFuzzingInput.Body)
			if err != nil {
				t.Fatalf("Failed to unmarshal fuzzy input: %s", err.Error())
			}

			err = json.Unmarshal(resp.Body.Bytes(), &jsonData.Body)
			if err != nil {
				t.Fatalf("Failed to unmarshal response body: %s", err.Error())
			}

			if jsonData.Body.Message != fmt.Sprintf("Hello post, %s!", unmarshalledFuzzingInput.Body.Name) {
				t.Fatalf("Unexpected response: %s", resp.Body.String())
			}
		},
	)
}
