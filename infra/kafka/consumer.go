package kafka

import (
	"context"

	"github.com/EdlanioJ/kbu-store/domain"
	kafka "github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	Reader          *kafka.Reader
	CategoryUsecase domain.CategoryUsecase
}

func NewKafkaConsumer(kafkaURLs []string, groupID string) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     kafkaURLs,
		GroupID:     groupID,
		GroupTopics: []string{"store.catetory.create", "store.catetory.update"},
		MinBytes:    10e3,
		MaxBytes:    10e6,
	})
	return &KafkaConsumer{
		Reader: reader,
	}
}

func (k *KafkaConsumer) Consume() {
	defer k.Reader.Close()
	log.Info("\u001b[92mStart Consuming...\u001b[0m")
	for {
		m, err := k.Reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}

		k.processMessage(m)
	}

}

func (k *KafkaConsumer) processMessage(msg kafka.Message) {
	ctx := context.Background()

	switch topic := msg.Topic; topic {
	case "store.catetory.create":
		err := k.createCategory(ctx, msg.Value)
		if err != nil {
			log.WithFields(log.Fields{
				"topic": topic,
				"msg":   string(msg.Value),
			}).Error(err)
		}
	case "store.catetory.update":
		err := k.updateCategory(ctx, msg.Value)
		if err != nil {
			log.WithFields(log.Fields{
				"topic": topic,
				"msg":   string(msg.Value),
			}).Error(err)
		}
	default:
		log.WithField("topic", topic).Warn("Invalid msg: ", string(msg.Value))
	}
}

func (k *KafkaConsumer) createCategory(ctx context.Context, data []byte) error {
	category := new(domain.Category)
	err := category.ParseJson(data)
	if err != nil {
		return err
	}
	return k.CategoryUsecase.Create(ctx, category)
}

func (k *KafkaConsumer) updateCategory(ctx context.Context, data []byte) error {
	category := new(domain.Category)
	err := category.ParseJson(data)
	if err != nil {
		return err
	}
	return k.CategoryUsecase.Update(ctx, category)
}
