package greeting

import (
	"context"
	"fmt"
)

func Get(ctx context.Context, input *InputGreeting) (resp *OutputGreeting, err error) {
	resp = &OutputGreeting{}
	resp.Body.Message = fmt.Sprintf("Hello get, %s!", input.Name)

	return
}

func Post(ctx context.Context, input *PostInputGreeting) (resp *OutputGreeting, err error) {
	resp = &OutputGreeting{}
	resp.Body.Message = fmt.Sprintf("Hello post, %s!", input.Body.Name)

	return resp, nil
}
