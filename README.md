# Features
This aims to be a production ready template repo that uses Go, Huma, Gin, Gorm, Postgres, Docker, and OTel.

* Example reads + write api endpoints with both ORM and raw SQL to Postgres.
* Test coverage for the routes including one fuzzy test example.
* Automatic DB migrations when running locally so the GORM models can be updated and the changes reflected in postgres.
* Runs 100% in Docker using Docker Compose and a simple make file with common commands.
* Logging. Using the ZeroLog framework all logs will go to both to Stdout and a log file.
* OTel telemetry + metrics + logging. OTel collector plus 3 backends (Jaeger, Zipkin, Prometheus)
* Sets and reads from Env vars in both docker and locally.
* OpenAPI spec generated and enforced at run time (Huma). Think FastAPI.
* Linting. Has most of the GolangCi-lint linters enables.
* Tests + linting enforced by Github actions on pullrequest and merge queue
* Builds and pushes a production image using Google distroless base image to Githubs container artafact storage.
* 37MB production docker image.


# Getting started
* Create a new repo using `Use this as template` in the Github UI.
* Find and replace `REPLACEME` with the name of your project.
* Find and replace `github.com/Sourceware-Lab/go-huma-gin-postgres-template` with the path of your repo
* Copy the `example.env` file to `.env`.
* Run `make run`

## Local go
If you don't want the backend to run in
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

# OTel
https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/examples/demo/docker-compose.yaml
https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/github.com/gin-gonic/gin/otelgin/example/server.go
* Jaeger at http://0.0.0.0:16686
* Zipkin at http://0.0.0.0:9411
* Prometheus at http://0.0.0.0:9090


# Additional Reading
Checkout the [docs](docs/index.md) dir. It contains files with additional information.


# OTEL
https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/examples/demo/docker-compose.yaml
https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/github.com/gin-gonic/gin/otelgin/example/server.go
* Jaeger at http://0.0.0.0:16686
* Zipkin at http://0.0.0.0:9411
* Prometheus at http://0.0.0.0:9090
