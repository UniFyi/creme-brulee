package config

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GetEnvDuration(envName string) (time.Duration, error) {
	val, err := GetEnv(envName)
	if err != nil {
		return 0, err
	}

	const timeFormat = `^([0-9]+)([sSmMhHdD])$`
	timeRegex := regexp.MustCompile(timeFormat)
	submatch := timeRegex.FindAllStringSubmatch(val, -1)
	if len(submatch) != 1 || len(submatch[0]) != 3 {
		return 0, fmt.Errorf("env duration is not in correct format %v", timeFormat)
	}
	amountStr := submatch[0][1]
	unit := strings.ToLower(submatch[0][2])

	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		return 0, fmt.Errorf("critical, problem in source code, regex should've ensured that prefix is int %v", val)
	}
	d := time.Duration(amount)

	switch unit {
	case "s":
		return d * time.Second, nil
	case "m":
		return d * time.Minute, nil
	case "h":
		return d * time.Hour, nil
	case "d":
		return d * 24 * time.Hour, nil
	}

	return 0, fmt.Errorf("critical, problem in source code, regex should've ensured that suffix is supported %v", val)
}