package kmanager

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/unifyi/creme-brulee/config"
	"net/http"
)

type InternalAPIServer struct {
	ctx           context.Context
	router        *gin.Engine
	healthChecker *KafkaHealthChecker
}

func NewKafkaHealthCheckServer(ctx context.Context, cfg *config.KafkaConfig, ginLogLevel string) (*InternalAPIServer, error) {
	gin.SetMode(ginLogLevel)
	router := gin.New()

	hChecker, err := newHealthzChecker(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &InternalAPIServer{
		ctx:    ctx,
		router: router,
		healthChecker: hChecker,
	}, nil
}

func (s *InternalAPIServer) Start() error {
	log := ctxlogrus.Extract(s.ctx)
	log.Info("starting serving internal APIs")

	s.initializeRoutes()

	defer s.healthChecker.Cleanup()
	return s.router.Run(":3000")
}

func (s *InternalAPIServer) initializeRoutes() {

	checkRoutes := s.router.Group("/checks")
	{
		checkRoutes.GET("/healthz", healthz(s.healthChecker))
	}
}

func healthz(hc *KafkaHealthChecker) gin.HandlerFunc {

	return func(c *gin.Context) {
		if hc.IsHealthy() {
			c.JSON(http.StatusNoContent, nil)
			return
		}

		c.JSON(http.StatusServiceUnavailable, nil)
		return
	}

}
