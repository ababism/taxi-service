package kafkaproducer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/juju/zaputil/zapctx"
	"github.com/segmentio/kafka-go"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/service/adapters"
	global "go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var _ adapters.KafkaClient = &KafkaProducer{}

type KafkaProducer struct {
	producer *kafka.Writer
}

func NewKafkaProducer(cfg *Config) *KafkaProducer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  cfg.Brokers,
		Topic:    cfg.Topic,
		Balancer: &kafka.LeastBytes{},
	})

	return &KafkaProducer{
		producer: writer,
	}
}

func (kp *KafkaProducer) SendTripStatusUpdate(ctx context.Context, trip domain.Trip, commandType domain.CommandType, reason *string) error {
	log := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	newCtx, span := tr.Start(ctx, "driver.repository.mongo: GetTripsByStatus")
	defer span.End()

	tc := ToTripCommand(trip, commandType, reason)

	message, err := json.Marshal(tc)
	if err != nil {
		log.Error("failed to marshal TripCommand to message:", zap.Error(err))
		return fmt.Errorf("failed to marshal TripCommand to message: %w", domain.ErrInternal)
	}

	err = kp.SendMessageToKafka(newCtx, message)
	if err != nil {
		log.Error("failed write message to kafka:", zap.Error(err))
		return fmt.Errorf("failed to send message to Kafka: %w", domain.ErrInternal)
	}
	return nil
}

func (kp *KafkaProducer) SendMessageToKafka(ctx context.Context, message []byte) error {
	return kp.producer.WriteMessages(ctx, kafka.Message{
		Value: message,
	})
}
