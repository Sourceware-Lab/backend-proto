package greeting

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"

	oteltrace "go.opentelemetry.io/otel/trace"

	"github.com/Sourceware-Lab/go-huma-gin-postgres-template/config"
)

func Get(c context.Context, input *InputGreeting) (*OutputGreeting, error) {
	resp := &OutputGreeting{}
	_, span := config.Tracer.Start(c, "getUser", oteltrace.WithAttributes(attribute.String("name", input.Name)))

	defer span.End()

	resp.Body.Message = fmt.Sprintf("Hello get, %s!", input.Name)

	return resp, nil
}

func Post(_ context.Context, input *PostInputGreeting) (*OutputGreeting, error) {
	resp := &OutputGreeting{}
	resp.Body.Message = fmt.Sprintf("Hello post, %s!", input.Body.Name)

	return resp, nil
}
