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
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"time"
)

const (
	acceptCommandType string = "trip.event.accept"
	cancelCommandType string = "trip.event.cancel"
	endCommandType    string = "trip.event.end"
	startCommandType  string = "trip.event.start"
	createCommandType string = "trip.event.create"
)

// В идеале запускать отдельным процессом в докере, но мы не успели красоту навести. :(

type KafkaConsumer struct {
	reader        *kafka.Reader
	driverService adapters.DriverService
}

func NewKafkaConsumer(cfg *Config, driverService adapters.DriverService) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Brokers,
		Topic:          cfg.Topic,
		GroupID:        cfg.IdGroup,
		MinBytes:       cfg.MinBytes,
		MaxBytes:       cfg.MaxBytes,
		SessionTimeout: 100 * time.Second,
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
			logger.Error("error while reading message from kafka", zap.Error(err))
			continue
		}

		tr := global.Tracer(domain.ServiceName)
		ctxTrace, span := tr.Start(ctx, "driver.daemon.kafkaConsumer: ConsumeMessage", trace.WithNewRoot())
		ctxLog := zapctx.WithLogger(ctxTrace, logger)
		var event CreatedTripEvent
		errCreate := json.Unmarshal(message.Value, &event)
		if errCreate != nil {
			logger.Debug("error unmarshalling message from kafka (unsupported schema)", zap.Error(err), zap.ByteString("json:", message.Value))
			continue
		}
		if errCreate == nil && event.Type == createCommandType {
			// TRIP_CREATED
			logger.Debug(fmt.Sprintf("Event id=%s caught with type=%s, created at=%s", event.ID, event.Type, event.Time), zap.Error(err))
			var ctEvent CreatedTripEvent
			ctEvent = event
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
			// close span and to next message in cycle
			span.End()
			continue
		}

		// events updating status
		var suEvent StatusUpdateEvent
		err = json.Unmarshal(message.Value, &suEvent)
		if err != nil {
			logger.Error("error while unmarshalling message value from kafka to StatusUpdateEvent", zap.Error(err))
			continue
		}
		tripId, err := ParseUUID(suEvent.Data.TripID)
		if err != nil {
			logger.Error(fmt.Sprintf("error while parsing TripId=%s to UUID", suEvent.RequestId), zap.Error(err))
			continue
		}
		switch event.Type {
		case acceptCommandType:
			err = kc.driverService.UpdateTripStatusAndDriver(ctxLog, tripId, suEvent.Data.DriverId, domain.TripStatuses.GetDriverFound())
		case cancelCommandType:
			err = kc.driverService.UpdateTripStatus(ctxLog, tripId, domain.TripStatuses.GetCanceled())
		case endCommandType:
			err = kc.driverService.UpdateTripStatus(ctxLog, tripId, domain.TripStatuses.GetEnded())
		case startCommandType:
			err = kc.driverService.UpdateTripStatus(ctxLog, tripId, domain.TripStatuses.GetStarted())
		default:
			logger.Debug(fmt.Sprintf("unkown kafka message type=%s", suEvent.Type), zap.Error(err))
			continue
		}
		if err != nil {
			logger.Error(fmt.Sprintf("error change status for %s of trip id=%s", suEvent.Type, tripId), zap.Error(err))
			continue
		}
		span.End()
	}
}
