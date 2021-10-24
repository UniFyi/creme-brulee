package logging

import (
	"context"
	"github.com/UniFyi/creme-brulee/pkg/config"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
)

func EnhanceContextWithLogger(ctx context.Context, cfg *config.BaseConfig) context.Context {
	logger := logrus.StandardLogger()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetReportCaller(true)

	if level, err := logrus.ParseLevel(cfg.LogLevel); err != nil {
		logger.Errorf("invalid log level %v", cfg.LogLevel)
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(level)
	}

	return ctxlogrus.ToContext(ctx, logrus.NewEntry(logger))
}
