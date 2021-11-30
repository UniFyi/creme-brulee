package config

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaConfig struct {
	Host     string
	Username string
	Password string
	SslMode  string
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
		Host:     envKafkaHost,
		Username: envKafkaUsername,
		Password: envKafkaPassword,
		SslMode:  envKafkaSslMode,
	}, nil
}

func (cfg *KafkaConfig) GetKafkaConfigMap(consumerGroup string) *kafka.ConfigMap {
	if cfg.SslMode != "disable" {
		return &kafka.ConfigMap{
			"bootstrap.servers":  cfg.Host,
			"group.id":           consumerGroup,
			"enable.auto.commit": false,
			"security.protocol":  "SASL_SSL",
			"sasl.mechanisms":    "PLAIN",
			"sasl.username":      cfg.Username,
			"sasl.password":      cfg.Password,
		}
	}

	return &kafka.ConfigMap{
		"bootstrap.servers":  cfg.Host,
		"group.id":           consumerGroup,
		"enable.auto.commit": false,
	}
}
