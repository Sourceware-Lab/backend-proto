package healthcheck

import (
	"context"
)

func Get(ctx context.Context, input *InputHealthcheck) (*OutputHealthcheck, error) {
	resp := &OutputHealthcheck{}
	resp.Status = 200

	return resp, nil
}
