package greeting

import (
	"context"
	"fmt"
)

func Get(ctx context.Context, input *InputGreeting) (*OutputGreeting, error) {
	resp := &OutputGreeting{}
	resp.Body.Message = fmt.Sprintf("Hello get, %s!", input.Name)

	return resp, nil
}

func Post(ctx context.Context, input *PostInputGreeting) (*OutputGreeting, error) {
	resp := &OutputGreeting{}
	resp.Body.Message = fmt.Sprintf("Hello post, %s!", input.Body.Name)

	return resp, nil
}
