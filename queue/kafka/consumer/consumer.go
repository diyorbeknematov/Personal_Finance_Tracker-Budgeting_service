package consumer

import (
	"context"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer interface {
	ConsumeMessages(ctx context.Context, handler func(message []byte)) error
	Close() error
}

type kafkaConsumerImpl struct {
	reader *kafka.Reader
	logger *slog.Logger
}

func NewKafkaConsumer(brokers []string, topic string, groupId string, logger *slog.Logger) KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupId,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	logger = logger.With("topic", topic, "group_id", groupId)

	return &kafkaConsumerImpl{
		reader: reader,
		logger: logger,
	}
}

func (k *kafkaConsumerImpl) ConsumeMessages(ctx context.Context, handler func(message []byte)) error {
	defer k.Close()

	for {
		select {
		case <-ctx.Done():
			k.logger.Info("Consumer shutting down")
			return nil
		default:
			m, err := k.reader.ReadMessage(ctx)
			if err != nil {
				k.logger.Error("Error reading message", "error", err)
				return err
			}

			go handler(m.Value) 
		}
	}
}

func (k *kafkaConsumerImpl) Close() error {
	return k.reader.Close()
}
