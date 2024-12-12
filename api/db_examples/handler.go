package greeting

import (
	"context"
	"fmt"
)

func PostRawSql(ctx context.Context, input *PostBodyInput) (*Output, error) {
	resp := &Output{}
	resp.Body.Message = fmt.Sprintf("Hello post, %s!", input.Body.Name)
	return resp, nil
}

func PostOrm(ctx context.Context, input *PostBodyInput) (*Output, error) {
	resp := &Output{}
	resp.Body.Message = fmt.Sprintf("Hello post, %s!", input.Body.Name)
	return resp, nil
}
