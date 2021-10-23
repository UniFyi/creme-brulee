package devsetup

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/unifyi/creme-brulee/pkg/config"
)

func CreateKafkaTopics(ctx context.Context, baseConf *config.BaseConfig, cfg *config.KafkaConfig, topicNames []string) {
	log := ctxlogrus.Extract(ctx)

	if baseConf.Env != "testing" {
		return
	}

	ka, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": cfg.KafkaHost,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to init kafka admin client %v", err))
	}

	topics := make([]kafka.TopicSpecification, len(topicNames))
	for i, name := range topicNames {
		topics[i] = kafka.TopicSpecification{
			Topic:             name,
			NumPartitions:     1,
			ReplicationFactor: 1,
		}
	}

	_, err = ka.CreateTopics(ctx, topics)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to auto create topics %v", err))
	}
}
