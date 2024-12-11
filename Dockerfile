FROM golang:1.23.4@sha256:574185e5c6b9d09873f455a7c205ea0514bfd99738c5dc7750196403a44ed4b7  AS base

ENV GOFLAGS="-buildvcs=false"

WORKDIR /app

FROM base AS build

COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/backend

FROM build AS final

COPY --from=build --chown=777 /app/backend /app/backend

EXPOSE 4000:4000

CMD ["./backend"]