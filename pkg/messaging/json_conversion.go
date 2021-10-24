package messaging

import (
	"encoding/json"
	"fmt"
)

func ConvertToJson(data interface{}) (string, error) {
	result, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to convert %T to json", data)
	}
	return string(result), nil
}
