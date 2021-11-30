package gintonic

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/unifyi/creme-brulee/messaging"
	"net/http"
)

type EndpointHandler func(c *gin.Context, ctx context.Context, uri map[string]uuid.UUID, data EndpointData)
type DataProvider func() interface{}
type StepFunc func(c *gin.Context, incoming *PipedData) *PipedData

type EndpointData struct {
	Payload     interface{}
	QueryParams interface{}
}

type endpointBuilder struct {
	rootCtx         context.Context
	orderedHandlers []StepFunc
}

type PipedData struct {
	completed bool
	ctx       context.Context
	uri       map[string]uuid.UUID
	data      EndpointData
}

func NewEndpointBuilder(ctx context.Context) *endpointBuilder {
	return &endpointBuilder{
		rootCtx:         ctx,
		orderedHandlers: make([]StepFunc, 0),
	}
}

func (eb *endpointBuilder) UserScoped() *endpointBuilder {
	eb.orderedHandlers = append(eb.orderedHandlers, func(c *gin.Context, incoming *PipedData) *PipedData {
		log := ctxlogrus.Extract(incoming.ctx)
		userID := c.GetHeader("UniFyi-User-Id")
		if userID == "" {
			log.Errorf("request is missing header 'UniFyi-User-Id'")
			c.JSON(http.StatusInternalServerError, nil)
			incoming.completed = true
			return incoming
		}
		userUUID, err := uuid.Parse(userID)
		if err != nil {
			log.Debugf("header 'UniFyi-User-Id' is not in uuid format %v", userID)
			c.JSON(http.StatusInternalServerError, nil)
			incoming.completed = true
			return incoming
		}

		incoming.ctx = context.WithValue(incoming.ctx, "userID", userUUID)
		incoming.completed = false
		return incoming
	})
	return eb
}

func (eb *endpointBuilder) WithURI(uriIdentifiers ...string) *endpointBuilder {
	eb.orderedHandlers = append(eb.orderedHandlers, func(c *gin.Context, incoming *PipedData) *PipedData {
		for _, uriID := range uriIdentifiers {
			identifier := c.Param(uriID)
			parsedUUID, err := uuid.Parse(identifier)
			if err != nil {
				c.JSON(http.StatusBadRequest, messaging.CreateInvalidFieldError(uriID, "not uuid format"))
				incoming.completed = true
				return incoming
			}

			incoming.uri[uriID] = parsedUUID
		}
		incoming.completed = false
		return incoming
	})
	return eb
}

func (eb *endpointBuilder) WithQueryParams(queryParamsProvider DataProvider) *endpointBuilder {
	eb.orderedHandlers = append(eb.orderedHandlers, func(c *gin.Context, incoming *PipedData) *PipedData {
		queryParams := queryParamsProvider()
		err := c.Bind(queryParams)
		if err != nil {
			c.JSON(http.StatusBadRequest, messaging.CreateGenericError("invalid query params"))
			incoming.completed = true
			return incoming
		}
		incoming.data.QueryParams = queryParams
		incoming.completed = false
		return incoming
	})
	return eb
}

func (eb *endpointBuilder) WithPayload(payloadProvider DataProvider) *endpointBuilder {
	eb.orderedHandlers = append(eb.orderedHandlers, func(c *gin.Context, incoming *PipedData) *PipedData {
		payload := payloadProvider()
		err := c.BindJSON(payload)
		if err != nil {
			c.JSON(http.StatusBadRequest, messaging.CreateGenericError("invalid payload"))
			incoming.completed = true
			return incoming
		}
		incoming.data.Payload = payload
		incoming.completed = false
		return incoming
	})
	return eb
}

func (eb *endpointBuilder) BuildFrom(endpointHandler EndpointHandler) gin.HandlerFunc {
	return func(c *gin.Context) {

		pipedData := &PipedData{
			completed: false,
			ctx:       eb.rootCtx,
			uri:       make(map[string]uuid.UUID, 0),
		}
		for _, handler := range eb.orderedHandlers {
			pipedData = handler(c, pipedData)
			if pipedData.completed {
				return
			}
		}

		endpointHandler(c, pipedData.ctx, pipedData.uri, pipedData.data)
	}
}
