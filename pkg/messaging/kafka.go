package messaging

import (
	"context"
	"github.com/UniFyi/creme-brulee/pkg/config"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"gorm.io/gorm"
)

type MessageConsumer struct {
	consumer   *kafka.Consumer
	topicNames []string
	db         *gorm.DB
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
