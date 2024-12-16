# Getting started
* Find and replace `REPLACEME` with the name of your project
* Copy the `example.env` file to `.env`
* Set any required env var values

## Local go
* https://go.dev/doc/install
* `go install github.com/air-verse/air@latest`
* run `air`

## Local docker
* https://docs.docker.com/compose/install/
* run `make run`

# Dependencies
## Adding
run `go get <url for go module>`

## Updating
Run `go get -u ./...`

# Additional Reading
Checkout the [docs](docs/index.md) dir. It contains files with additional information.


# OTEL
https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/examples/demo/docker-compose.yaml
https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/github.com/gin-gonic/gin/otelgin/example/server.go
* Jaeger at http://0.0.0.0:16686
* Zipkin at http://0.0.0.0:9411
* Prometheus at http://0.0.0.0:9090
