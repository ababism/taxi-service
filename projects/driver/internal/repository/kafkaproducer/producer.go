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
	writer := kafka.Writer{
		Addr:     kafka.TCP(cfg.Broker),
		Topic:    cfg.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaProducer{
		producer: &writer,
	}
}

func (kp *KafkaProducer) SendTripStatusUpdate(ctx context.Context, trip domain.Trip, commandType domain.CommandType, reason *string) error {
	logger := zapctx.Logger(ctx)
	logger.Debug("[kafka_producer] config:", zap.String("topic", kp.producer.Topic))

	tr := global.Tracer(domain.ServiceName)
	newCtx, span := tr.Start(ctx, "driver.repository: SendTripStatusUpdate")
	defer span.End()

	tc := ToTripCommand(trip, commandType, reason)

	message, err := json.Marshal(tc)
	if err != nil {
		logger.Error("failed to marshal TripCommand to message:", zap.Error(err))
		return fmt.Errorf("failed to marshal TripCommand to message: %w", domain.ErrInternal)
	}

	err = kp.SendMessageWithKaKafka(newCtx, message)
	if err != nil {
		logger.Error("failed write message to kafka:", zap.Error(err))
		return fmt.Errorf("failed to send message to Kafka: %w", domain.ErrInternal)
	}
	return nil
}

func (kp *KafkaProducer) SendMessageWithKaKafka(ctx context.Context, message []byte) error {
	return kp.producer.WriteMessages(ctx, kafka.Message{
		Value: message,
	})
}
