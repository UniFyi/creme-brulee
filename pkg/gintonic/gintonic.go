package gintonic

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"net/http"
)

type EndpointHandler func(c *gin.Context, ctx context.Context, uri map[string]uuid.UUID, data EndpointData)
type DataProvider func() interface{}

type EndpointData struct {
	Payload     interface{}
	QueryParams interface{}
}

type endpointBuilder struct {
	ctx             context.Context
	orderedHandlers []gin.HandlerFunc
	completed       bool
	uri             map[string]uuid.UUID
	data            EndpointData
}

func NewEndpointBuilder(ctx context.Context) *endpointBuilder {
	return &endpointBuilder{
		ctx:             ctx,
		orderedHandlers: make([]gin.HandlerFunc, 0),
		uri:             make(map[string]uuid.UUID, 0),

	}
}

func (eb *endpointBuilder) UserScoped() *endpointBuilder {
	eb.orderedHandlers = append(eb.orderedHandlers, func(c *gin.Context) {
		if eb.completed {
			return
		}

		log := ctxlogrus.Extract(eb.ctx)
		userID := c.GetHeader("UniFyi-User-Id")
		if userID == "" {
			log.Errorf("request is missing header 'UniFyi-User-Id'")
			c.JSON(http.StatusInternalServerError, nil)
			eb.completed = true
			return
		}
		userUUID, err := uuid.Parse(userID)
		if err != nil {
			log.Debugf("header 'UniFyi-User-Id' is not in uuid format %v", userID)
			c.JSON(http.StatusInternalServerError, nil)
			eb.completed = true
			return
		}

		eb.ctx = context.WithValue(eb.ctx, "userID", userUUID)
	})
	return eb
}

func (eb *endpointBuilder) WithURI(uriIdentifiers ...string) *endpointBuilder {
	eb.orderedHandlers = append(eb.orderedHandlers, func(c *gin.Context) {
		if eb.completed {
			return
		}

		log := ctxlogrus.Extract(eb.ctx)
		for _, uriID := range uriIdentifiers {
			identifier := c.Param(uriID)
			parsedUUID, err := uuid.Parse(identifier)
			if err != nil {
				log.Debug("%v id is not uuid", uriID)
				c.JSON(http.StatusBadRequest, nil)
				eb.completed = true
				return
			}

			eb.uri[uriID] = parsedUUID
		}
	})
	return eb
}

func (eb *endpointBuilder) WithQueryParams(queryParamsProvider DataProvider) *endpointBuilder {
	eb.orderedHandlers = append(eb.orderedHandlers, func(c *gin.Context) {
		if eb.completed {
			return
		}

		queryParams := queryParamsProvider()
		log := ctxlogrus.Extract(eb.ctx)
		err := c.Bind(queryParams)
		if err != nil {
			log.Debugf("invalid query params")
			c.JSON(http.StatusBadRequest, nil)
			eb.completed = true
			return
		}
		eb.data.QueryParams = queryParams
	})
	return eb
}

func (eb *endpointBuilder) WithPayload(payloadProvider DataProvider) *endpointBuilder {
	eb.orderedHandlers = append(eb.orderedHandlers, func(c *gin.Context) {
		if eb.completed {
			return
		}

		payload := payloadProvider()
		log := ctxlogrus.Extract(eb.ctx)
		err := c.BindJSON(payload)
		if err != nil {
			log.Debugf("invalid payload")
			c.JSON(http.StatusBadRequest, nil)
			eb.completed = true
			return
		}
		eb.data.Payload = payload
	})
	return eb
}

func (eb *endpointBuilder) BuildFrom(endpointHandler EndpointHandler) gin.HandlerFunc {
	return func(c *gin.Context) {

		for _, handler := range eb.orderedHandlers {
			handler(c)
		}
		if eb.completed {
			return
		}

		endpointHandler(c, eb.ctx, eb.uri, eb.data)
	}
}
