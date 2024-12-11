package greeting

import (
	"context"
	"fmt"
)

func Get(ctx context.Context, input *Input) (*Output, error) {
	resp := &Output{}
	resp.Body.Message = fmt.Sprintf("Hello get , %s!", input.Name)
	return resp, nil
}

func Post(ctx context.Context, input *Input) (*Output, error) {
	resp := &Output{}
	resp.Body.Message = fmt.Sprintf("Hello post , %s!", input.Name)
	return resp, nil
}
