package kmanager

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/unifyi/creme-brulee/config"
	"gorm.io/gorm"
)

const (
	HealthTopic          = "Healthz"
	HealthzConsumerGroup = "healthz-consumer-group"
)

type MessageConsumer struct {
	consumer   *kafka.Consumer
	topicNames []string
	db         *gorm.DB
}

type KafkaHealthChecker struct {
	consumer *kafka.Consumer
}

func NewMessageConsumer(ctx context.Context, db *gorm.DB, cfg *config.KafkaConfig, consumerGroup string, topicNames []string) *MessageConsumer {
	log := ctxlogrus.Extract(ctx)

	kc, kafkaError := kafka.NewConsumer(cfg.GetKafkaConfigMap(consumerGroup))
	if kafkaError != nil {
		log.Fatal(kafkaError)
	}
	return &MessageConsumer{
		consumer:   kc,
		topicNames: topicNames,
		db:         db,
	}
}

type TopicHandler func(ctx context.Context, db *gorm.DB, msg *kafka.Message) error

func (mc *MessageConsumer) Start(ctx context.Context, handleMessage TopicHandler) error {
	log := ctxlogrus.Extract(ctx)
	log.Info("starting kafka consumer")

	err := mc.consumer.SubscribeTopics(mc.topicNames, nil)
	if err != nil {
		log.Errorf("failed to subscirbe to kafka topics %v", err)
		return err
	}
	defer mc.consumer.Close()

	for {
		// TODO how often do we ask kafka for messages?
		if msg, err := mc.consumer.ReadMessage(-1); err == nil {
			if len(msg.Key) != 0 { // ignore heartbeat messages
				if err = handleMessage(ctx, mc.db, msg); err != nil {
					log.Error("message will be NOT committed")
					// TODO if we couldn't consume event from kafka we need to retry
					// TODO do we fail the pod since it couldn't consume event?
					continue
				}
			}

			if _, err := mc.consumer.CommitMessage(msg); err != nil {
				log.Warnf("couldn't commit message %v", err)
			}
		} else {
			log.Warnf("consumer kafka error: %v (%v)\n", err, msg)
		}
	}
}

func newHealthzChecker(ctx context.Context, cfg *config.KafkaConfig) (*KafkaHealthChecker, error) {
	log := ctxlogrus.Extract(ctx)
	log.Info("starting healthz kafka consumer")

	kc, kafkaError := kafka.NewConsumer(cfg.GetKafkaConfigMap(HealthzConsumerGroup))
	if kafkaError != nil {
		log.Fatal(kafkaError)
	}
	err := kc.SubscribeTopics([]string{HealthTopic}, nil)
	if err != nil {
		log.Errorf("failed to subscirbe to kafka topics %v", err)
		return nil, err
	}

	return &KafkaHealthChecker{
		consumer: kc,
	}, nil
}

func (h *KafkaHealthChecker) Cleanup() error {
	return h.consumer.Close()
}

func (h *KafkaHealthChecker) IsHealthy() bool {
	_, err := h.consumer.ReadMessage(-1)
	// is healthy only if error is null
	return err == nil
}
