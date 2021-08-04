package kafka

import (
	"context"

	"github.com/EdlanioJ/kbu-store/app/config"
	"github.com/opentracing/opentracing-go"
	kafka "github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

type KafkaProducer struct {
	Writer *kafka.Writer
}

func NewKafkaProducer(cfg *config.Config) *KafkaProducer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:     cfg.Kafka.Brokers,
		Balancer:    &kafka.LeastBytes{},
		Logger:      kafka.LoggerFunc(log.Debugf),
		ErrorLogger: kafka.LoggerFunc(log.Errorf),
	})

	return &KafkaProducer{
		Writer: writer,
	}
}

func (k *KafkaProducer) Publish(ctx context.Context, msg, topic string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "kafkaProducer.Publish")
	defer span.Finish()

	message := kafka.Message{
		Topic: topic,
		Value: []byte(msg),
	}

	return k.Writer.WriteMessages(ctx, message)
}

func (k *KafkaProducer) Close() {
	_ = k.Writer.Close()
}
