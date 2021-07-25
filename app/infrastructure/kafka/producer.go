package kafka

import (
	"context"

	"github.com/EdlanioJ/kbu-store/app/config"
	kafka "github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	Writer *kafka.Writer
}

func NewKafkaProducer(cfg *config.Config) *KafkaProducer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  cfg.Kafka.Brokers,
		Balancer: &kafka.LeastBytes{},
	})

	return &KafkaProducer{
		Writer: writer,
	}
}

func (k *KafkaProducer) Publish(ctx context.Context, msg, topic string) error {
	message := kafka.Message{
		Topic: topic,
		Value: []byte(msg),
	}

	return k.Writer.WriteMessages(ctx, message)
}

func (k *KafkaProducer) Close() {
	_ = k.Writer.Close()
}
