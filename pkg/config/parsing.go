package config

import (
	"fmt"
	"os"
	"strconv"
)

func GetEnvBool(envName string) (bool, error) {
	val, err := GetEnv(envName)
	if err != nil {
		return false, err
	}
	result, err := strconv.ParseBool(val)
	if err != nil {
		return false, err
	}
	return result, nil
}

func GetEnv(envName string) (string, error) {
	val, ok := os.LookupEnv(envName)
	if ok {
		return val, nil
	}

	return "", fmt.Errorf("missing env var: %s", envName)
}
