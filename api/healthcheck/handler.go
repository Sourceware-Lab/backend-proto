package healthcheck

import (
	"context"
)

func Get(ctx context.Context, input *InputHealthcheck) (resp *OutputHealthcheck, err error) {
	resp = &OutputHealthcheck{}
	resp.Status = 200
	return
}
