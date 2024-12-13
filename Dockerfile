FROM golang:1.23.4@sha256:574185e5c6b9d09873f455a7c205ea0514bfd99738c5dc7750196403a44ed4b7  AS base

ENV GOFLAGS="-buildvcs=false"

ARG UID=1000
ARG GID=$UID
ARG USERNAME=nonroot
ENV WORKDIR=/app


RUN addgroup --gid $GID $USERNAME && \
    adduser --uid $UID --gid $GID --disabled-password --gecos "" $USERNAME

WORKDIR $WORKDIR

RUN chown -R $USERNAME $WORKDIR

USER $USERNAME

FROM base AS local

RUN go install github.com/air-verse/air@latest
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/backend

FROM base AS production

COPY --from=local --chown=777 /app/backend /app/backend

CMD ["./backend"]
