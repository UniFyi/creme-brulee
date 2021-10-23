package gintonic

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"net/http"
)

type EndpointHandler func(c *gin.Context, ctx context.Context, uri map[string]uuid.UUID) gin.HandlerFunc

type endpointBuilder struct {
	c         *gin.Context
	ctx       context.Context
	completed bool
	uri       map[string]uuid.UUID
}

func NewEndpointBuilder(c *gin.Context, ctx context.Context) *endpointBuilder {
	return &endpointBuilder{
		c:   c,
		ctx: ctx,
	}
}

func (eb *endpointBuilder) UserScoped() *endpointBuilder {
	if eb.completed {
		return eb
	}

	log := ctxlogrus.Extract(eb.ctx)
	userID := eb.c.GetHeader("UniFyi-User-Id")
	if userID == "" {
		log.Errorf("request is missing header 'UniFyi-User-Id'")
		eb.c.JSON(http.StatusInternalServerError, nil)
		eb.completed = true
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Debugf("header 'UniFyi-User-Id' is not in uuid format %v", userID)
		eb.c.JSON(http.StatusInternalServerError, nil)
		eb.completed = true
	}

	eb.ctx = context.WithValue(eb.ctx, "userID", userUUID)
	return eb
}

func (eb *endpointBuilder) WithURI(uriIdentifiers ...string) *endpointBuilder {
	if eb.completed {
		return eb
	}

	log := ctxlogrus.Extract(eb.ctx)
	for _, uriID := range uriIdentifiers {
		identifier := eb.c.Param(uriID)
		parsedUUID, err := uuid.Parse(identifier)
		if err != nil {
			log.Debug("%v id is not uuid", uriID)
			eb.c.JSON(http.StatusBadRequest, nil)
			eb.completed = true
			return eb
		}

		eb.uri[uriID] = parsedUUID
	}

	return eb
}

func (eb *endpointBuilder) WithQueryParams(queryParams interface{}) *endpointBuilder {
	if eb.completed {
		return eb
	}

	log := ctxlogrus.Extract(eb.ctx)
	err := eb.c.Bind(queryParams)
	if err != nil {
		log.Debugf("invalid query params")
		eb.c.JSON(http.StatusBadRequest, nil)
		eb.completed = true
	}

	return eb
}

func (eb *endpointBuilder) WithPayload(payload interface{}) *endpointBuilder {
	if eb.completed {
		return eb
	}

	log := ctxlogrus.Extract(eb.ctx)
	err := eb.c.BindJSON(payload)
	if err != nil {
		log.Debugf("invalid payload")
		eb.c.JSON(http.StatusBadRequest, nil)
		eb.completed = true
	}

	return eb
}

func (eb *endpointBuilder) BuildFrom(endpointHandler EndpointHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if eb.completed {
			return
		}
		endpointHandler(eb.c, eb.ctx, eb.uri)
	}
}
