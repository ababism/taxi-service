package kafkaconsumer

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
	_ "time"
)

// В идеале запускать отдельным процессом в докере, но мы не успели красоту навести. :(

type KafkaConsumer struct {
	reader        *kafka.Reader
	driverService adapters.DriverService
}

func NewKafkaConsumer(cfg *Config, driverService adapters.DriverService) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Brokers,
		Topic:    cfg.Topic,
		GroupID:  cfg.IdGroup,
		MinBytes: cfg.MinBytes, // Should be low
		MaxBytes: cfg.MaxBytes,
	})

	return &KafkaConsumer{
		reader:        reader,
		driverService: driverService,
	}
}

func (kc *KafkaConsumer) Start(ctx context.Context) func(kc *KafkaConsumer) error {
	go kc.consumeMessages(ctx)

	return closeKafka
}

func closeKafka(kc *KafkaConsumer) error {
	return kc.reader.Close()
}

func (kc *KafkaConsumer) consumeMessages(mainCtx context.Context) {
	logger := zapctx.Logger(mainCtx)

	for {
		// Block (wait)

		// Создаю новый контекст и логгер для отслеживания trace
		ctx := context.Background()
		message, err := kc.reader.ReadMessage(ctx)
		if err != nil {
			logger.Debug("error while reading message from kafka", zap.Error(err))
			continue
		}

		tr := global.Tracer(domain.ServiceName)
		ctxTrace, span := tr.Start(ctx, "driver.service: UpdateTripStatus")
		ctxLog := zapctx.WithLogger(ctxTrace, logger)

		var event Event
		err = json.Unmarshal(message.Value, &event)
		if err != nil {
			logger.Error("error while unmarshalling message from kafka", zap.Error(err))
			continue
		}

		logger.Debug(fmt.Sprintf("Event id=%s caught with type=%s, created at=%s", event.ID, event.Type, event.Time), zap.Error(err))

		switch event.Type {
		case "trip.event.created":
			var ctEvent CreatedTripEvent
			err = json.Unmarshal(message.Value, &ctEvent)
			if err != nil {
				logger.Error("error while unmarshalling message value from kafka to CreatedTripEvent", zap.Error(err))
				continue
			}
			trip, err := ToDomainTrip(ctEvent.Data)
			if err != nil {
				logger.Error(fmt.Sprintf("error while parsing TripId=%s to UUID", ctEvent.ID), zap.Error(err))
				continue
			}
			err = kc.driverService.InsertTrip(ctxLog, trip)
			if err != nil {
				logger.Error(fmt.Sprintf("error while insertiong trip id=%s", trip.Id), zap.Error(err))
				continue
			}
		case "trip.event.accepted", "trip.event.canceled", "trip.event.started", "trip.event.ended":
			var suEvent StatusUpdateEvent
			err = json.Unmarshal(message.Value, &suEvent)
			if err != nil {
				logger.Error("error while unmarshalling message value from kafka to StatusUpdateEvent", zap.Error(err))
				continue
			}
			tripId, err := ParseUUID(suEvent.ID)
			if err != nil {
				logger.Error(fmt.Sprintf("error while parsing TripId=%s to UUID", suEvent.ID), zap.Error(err))
				continue
			}
			switch event.Type {
			case "trip.event.accepted":
				err = kc.driverService.UpdateTripStatus(ctxLog, tripId, domain.TripStatuses.GetDriverFound())
			case "trip.event.canceled":
				err = kc.driverService.UpdateTripStatus(ctxLog, tripId, domain.TripStatuses.GetCanceled())
			case "trip.event.ended":
				err = kc.driverService.UpdateTripStatus(ctxLog, tripId, domain.TripStatuses.GetEnded())
			case "trip.event.started":
				err = kc.driverService.UpdateTripStatus(ctxLog, tripId, domain.TripStatuses.GetStarted())
			}
			if err != nil {
				logger.Error(fmt.Sprintf("error change status for %s of trip id=%s", suEvent.Type, tripId), zap.Error(err))
				continue
			}
		default:
			logger.Debug(fmt.Sprintf("Unknown event type %s for consumer. It will be skipped.", event.Type), zap.Error(err))
			continue
		}

		span.End()

		// Process the received message
		// Example: fmt.Printf("Received message: %s\n", string(message.Value))

	}
}
