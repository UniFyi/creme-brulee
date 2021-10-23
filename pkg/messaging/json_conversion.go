package messaging

import (
	"encoding/json"
	"fmt"
)

type Convertable struct {}

func (n Convertable) ToJSON() (string, error) {
	return convertToJson(n)
}

func convertToJson(data interface{}) (string, error) {
	result, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to convert %T to json", data)
	}
	return string(result), nil
}
