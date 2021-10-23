package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
)

type EventTemplate struct {
	Template interface{}
}

func NewEventFromTemplate(ctx context.Context, data []byte, eventTemplates map[string]EventTemplate) (interface{}, error) {
	log := ctxlogrus.Extract(ctx)

	var objMap map[string]interface{}
	err := json.Unmarshal(data, &objMap)
	if err != nil {
		return nil, err
	}

	body := string(data)
	log.Debugf("incoming data from kafka %v", body)
	if eventType, ok := objMap["type"]; ok {
		event, ok := eventTemplates[eventType.(string)]
		if !ok {
			return nil, fmt.Errorf("unknown event type %v", eventType)
		}
		return newEvent(data, eventType, body, event.Template)
	}
	return nil, fmt.Errorf("missing key [type] in %v", body)
}

func newEvent(data []byte, eventType interface{}, body string, template interface{}) (interface{}, error) {
	if err := json.Unmarshal(data, template); err != nil {
		return nil, fmt.Errorf(
			"failed to decode userEvent of type %v with value %v because %v", eventType, body, err)
	}
	return template, nil
}

// EXAMPLE
//func NewUser(ctx context.Context, data []byte) (interface{}, error) {
//	var user *User
//	return NewEventFromTemplate(ctx, data, map[string]EventTemplate{
//		UserCreatedEventKey: {
//			Template: user,
//		},
//		UserUpdatedEventKey: {
//			Template: user,
//		},
//	})
//}
