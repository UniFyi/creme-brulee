package config

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaConfig struct {
	KafkaHost     string
	KafkaUsername string
	KafkaPassword string
	KafkaSslMode  string
}

func NewKafkaConfig() (*KafkaConfig, error) {

	envKafkaHost, err := GetEnv("KAFKA_HOST")
	if err != nil {
		return nil, err
	}
	envKafkaSslMode, err := GetEnv("KAFKA_SSL")
	if err != nil {
		return nil, err
	}

	var envKafkaUsername string
	var envKafkaPassword string
	if envKafkaSslMode != "disable" {
		envKafkaUsername, err = GetEnv("KAFKA_USERNAME")
		if err != nil {
			return nil, err
		}
		envKafkaPassword, err = GetEnv("KAFKA_PASSWORD")
		if err != nil {
			return nil, err
		}
	}

	return &KafkaConfig{
		KafkaHost:     envKafkaHost,
		KafkaUsername: envKafkaUsername,
		KafkaPassword: envKafkaPassword,
		KafkaSslMode:  envKafkaSslMode,
	}, nil
}

func (cfg *KafkaConfig) GetKafkaConfigMap(consumerGroup string) *kafka.ConfigMap {
	if cfg.KafkaSslMode != "disable" {
		return &kafka.ConfigMap{
			"bootstrap.servers":  cfg.KafkaHost,
			"group.id":           consumerGroup,
			"enable.auto.commit": false,
			"security.protocol":  "SASL_SSL",
			"sasl.mechanisms":    "PLAIN",
			"sasl.username":      cfg.KafkaUsername,
			"sasl.password":      cfg.KafkaPassword,
		}
	}

	return &kafka.ConfigMap{
		"bootstrap.servers":  cfg.KafkaHost,
		"group.id":           "spicy-mango-group",
		"enable.auto.commit": false,
	}
}
