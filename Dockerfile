FROM golang:1.23.4-alpine@sha256:6c5c9590f169f77c8046e45c611d3b28fe477789acd8d3762d23d4744de69812 AS base

ENV GOFLAGS="-buildvcs=false"

ARG UID=1000
ARG GID=$UID
ARG USERNAME=nonroot
ENV WORKDIR=/app


RUN addgroup -g "$GID" "$USERNAME" && adduser -u "$UID" -G "$USERNAME" -D -g "" "$USERNAME"

WORKDIR $WORKDIR

RUN chown -R $USERNAME $WORKDIR

USER $USERNAME

FROM base AS install

COPY go.mod go.sum ./
RUN go mod download

FROM install AS local

USER root
RUN apk add --no-cache bash
USER $USERNAME

RUN go install github.com/air-verse/air@latest
COPY . .


FROM install AS build
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/backend

# docker pull gcr.io/distroless/static-debian12:nonroot
FROM gcr.io/distroless/static-debian12@sha256:6cd937e9155bdfd805d1b94e037f9d6a899603306030936a3b11680af0c2ed58 AS production

COPY --from=build /app/backend /backend
EXPOSE 8080

ENTRYPOINT ["/backend"]
