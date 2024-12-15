package healthcheck

import (
	"context"
)

func Get(_ context.Context, _ *InputHealthcheck) (*OutputHealthcheck, error) {
	resp := &OutputHealthcheck{}
	resp.Status = 200

	return resp, nil
}
