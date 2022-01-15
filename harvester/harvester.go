package harvester

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/unifyi/creme-brulee/config"
	"github.com/unifyi/creme-brulee/gintonic"
)

func Prompt(ctx context.Context) {
	envHarvesterHost, err := config.GetEnv("HARVESTER_HOST")
	if err != nil {
		log := ctxlogrus.Extract(ctx)
		log.Warn("Harvester prompt failed, missing HARVESTER_HOST env var")
		return
	}
	_, _ = gintonic.MakeRequest(ctx, "POST", fmt.Sprintf("%v/prompt", envHarvesterHost), nil)
}
