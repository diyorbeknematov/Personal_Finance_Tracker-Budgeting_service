package consumer

import (
	"budgeting-service/service"
	"context"
	"log"
	"log/slog"
)

type KafkaMethods interface {
	CreateTransaction(ctx context.Context, topic string)
	UpdateBudget(ctx context.Context, topic string)
}

type kafkaMethodsImpl struct {
	brokers          []string
	msgBrokerService service.MsgBrokerService
	logger           *slog.Logger
}

func NewKafkaMethods(brokers []string, msgBrokerService service.MsgBrokerService, logger *slog.Logger) KafkaMethods {
	return &kafkaMethodsImpl{
		brokers:          brokers,
		msgBrokerService: msgBrokerService,
		logger:           logger,
	}
}

func (km *kafkaMethodsImpl) CreateTransaction(ctx context.Context, topic string) {
	reader := NewKafkaConsumer(km.brokers, topic, "", km.logger)
	defer reader.Close()

	log.Println("Starting consumer for topic", topic)

	err := reader.ConsumeMessages(ctx, km.msgBrokerService.CreateTransaction)
	if err != nil {
		km.logger.Error("Error consuming messages", "error", err)
		log.Println("Error consuming messages", "error", err)
		return
	}
}

func (km *kafkaMethodsImpl) UpdateBudget(ctx context.Context, topic string) {
	reader := NewKafkaConsumer(km.brokers, topic, "", km.logger)
	defer reader.Close()

	log.Println("Starting consumer for topic", topic)

	err := reader.ConsumeMessages(ctx, km.msgBrokerService.UpdateBudget)
	if err != nil {
		km.logger.Error("Error consuming messages", "error", err)
		log.Println("Error consuming messages", "error", err)
		return
	}
}
