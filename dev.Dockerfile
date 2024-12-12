FROM golang:1.23.4@sha256:574185e5c6b9d09873f455a7c205ea0514bfd99738c5dc7750196403a44ed4b7  AS base

ENV GOFLAGS="-buildvcs=false"

WORKDIR /app

FROM base AS final
RUN go install github.com/air-verse/air@latest
COPY go.mod go.sum ./
RUN go mod download


CMD ["air", "-c", ".air.toml"]
