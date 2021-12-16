package config

import (
	"time"
)

func GetEnvDuration(envName string) (time.Duration, error) {
	val, err := GetEnv(envName)
	if err != nil {
		return 0, err
	}

	return time.ParseDuration(val)
}