package harvester

import (
	"context"
	"github.com/unifyi/creme-brulee/gintonic"
)

func Prompt(ctx context.Context) {
	_, _ = gintonic.MakeRequest(ctx, "POST", "http://localhost:4000/prompt", nil)
}
