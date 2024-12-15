package greeting

import (
	"context"
	"fmt"
)

func Get(_ context.Context, input *InputGreeting) (*OutputGreeting, error) {
	resp := &OutputGreeting{}
	resp.Body.Message = fmt.Sprintf("Hello get, %s!", input.Name)

	return resp, nil
}

func Post(_ context.Context, input *PostInputGreeting) (*OutputGreeting, error) {
	resp := &OutputGreeting{}
	resp.Body.Message = fmt.Sprintf("Hello post, %s!", input.Body.Name)

	return resp, nil
}
