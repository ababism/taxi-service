package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/juju/zaputil/zapctx"
	"github.com/pkg/errors"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/domain"

	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/service/adapters"
	global "go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var _ adapters.DriverService = &driverService{}

type driverService struct {
	r              adapters.DriverRepository
	locationClient adapters.LocationClient
	kafkaClient    adapters.KafkaClient
}

func NewDriverService(driverRepository adapters.DriverRepository, locationClient adapters.LocationClient, kafkaClient adapters.KafkaClient) adapters.DriverService {
	return &driverService{
		r:              driverRepository,
		locationClient: locationClient,
		kafkaClient:    kafkaClient,
	}
}

// GetTripsByStatus returns all trips with the given status
func (s *driverService) GetTripsByStatus(ctx context.Context, status domain.TripStatus) ([]domain.Trip, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	newCtx, span := tr.Start(ctx, "driver.service: GetTripsByStatus")
	defer span.End()

	trips, err := s.r.GetTripsByStatus(newCtx, status)
	if err != nil {
		logger.Error("driver-service: get trips by status from repository", zap.Error(err))
		return nil, err
	}

	return trips, nil
}

func (s *driverService) GetTripByID(ctx context.Context, driverId uuid.UUID, tripId uuid.UUID) (*domain.Trip, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	newCtx, span := tr.Start(ctx, "driver.service: GetTripByID")
	defer span.End()

	// err if trip driver != nil and driver != driverId
	trip, err := s.r.GetTripByID(newCtx, tripId)
	if err != nil {
		logger.Error("driver-service: get trip from repository")
		return nil, err
	}
	if trip.DriverId != nil && *trip.DriverId != driverId.String() {
		return nil, errors.Wrap(domain.ErrAccessDenied, "trip driver id does not match passed id")
	}
	return trip, err
}

// InsertTrip inserts a trip
func (s *driverService) InsertTrip(ctx context.Context, trip domain.Trip) error {
	log := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	newCtx, span := tr.Start(ctx, "driver.service: InsertTrip")
	defer span.End()

	err := s.r.InsertTrip(newCtx, trip)
	if err != nil {
		log.Error("driver-service: insert trip in repository", zap.Error(err))
		return err
	}

	return nil
}

// UpdateTrip updates all fields of the given trip in the service layer
func (s *driverService) UpdateTrip(ctx context.Context, tripId uuid.UUID, updatedTrip domain.Trip) error {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	newCtx, span := tr.Start(ctx, "driver.service: UpdateTrip")
	defer span.End()

	// Call the repository method to update the trip
	err := s.r.UpdateTrip(newCtx, tripId, updatedTrip)
	if err != nil {
		logger.Error("driver-service: update trip in repository", zap.Error(err))
		return err
	}

	return nil
}

func (s *driverService) AcceptTrip(ctx context.Context, driverId uuid.UUID, tripId uuid.UUID) error {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	newCtx, span := tr.Start(ctx, "driver.service: AcceptTrip")
	defer span.End()

	ctx = zapctx.WithLogger(newCtx, logger)

	trip, err := s.r.GetTripByID(newCtx, tripId)
	if err != nil || trip == nil {
		logger.Error("can't get trip from repository", zap.Error(err))
		return err
	}

	if trip.Status == nil || *trip.Status != domain.TripStatuses.GetDriverSearch() {
		return errors.Wrap(domain.ErrAccessDenied, "trip doesn't need driver")
	}
	dId := driverId.String()
	trip.DriverId = &dId
	err = s.kafkaClient.SendTripStatusUpdate(newCtx, *trip, domain.TripCommandAccept, nil)
	if err != nil {
		logger.Error("can't send accept trip to kafka:", zap.Error(err))
		return err
	}

	return nil
}

func (s *driverService) CancelTrip(ctx context.Context, driverId uuid.UUID, tripId uuid.UUID, reason *string) error {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	newCtx, span := tr.Start(ctx, "driver.service: CancelTrip")
	defer span.End()

	ctx = zapctx.WithLogger(newCtx, logger)

	trip, err := s.r.GetTripByID(newCtx, tripId)
	if err != nil || trip == nil {
		logger.Error("can't get trip from repository", zap.Error(err))
		return err
	}
	//if trip.Status == nil || *trip.Status == domain.TripStatuses.GetDriverSearch() {
	//	return errors.Wrap(domain.ErrAccessDenied, "trip hasn't connected with driver yet")
	//}
	if trip.DriverId != nil && *trip.DriverId != driverId.String() {
		return errors.Wrap(domain.ErrAccessDenied, "trip driver id does not match passed id")
	}
	dId := driverId.String()
	trip.DriverId = &dId
	err = s.kafkaClient.SendTripStatusUpdate(newCtx, *trip, domain.TripCommandCancel, reason)
	if err != nil {
		logger.Error("can't send cancel trip command to kafka:", zap.Error(err))
		return err
	}
	return nil
}

func (s *driverService) StartTrip(ctx context.Context, driverId uuid.UUID, tripId uuid.UUID) error {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	newCtx, span := tr.Start(ctx, "driver.service: StartTrip")
	defer span.End()

	ctx = zapctx.WithLogger(newCtx, logger)

	trip, err := s.r.GetTripByID(newCtx, tripId)
	if err != nil || trip == nil {
		logger.Error("can't get trip from repository", zap.Error(err))
		return err
	}
	//  ask about *trip.Status != s.tripStatuses.GetOnPosition()
	//if trip.Status == nil || *trip.Status != domain.TripStatuses.GetOnPosition() || *trip.Status != domain.TripStatuses.GetDriverFound() {
	//	return errors.Wrap(domain.ErrAccessDenied, "trip hasn't connected with driver yet")
	//}
	if trip.DriverId != nil && *trip.DriverId != driverId.String() {
		return errors.Wrap(domain.ErrAccessDenied, "trip driver id does not match passed id")
	}
	dId := driverId.String()
	trip.DriverId = &dId
	err = s.kafkaClient.SendTripStatusUpdate(newCtx, *trip, domain.TripCommandStart, nil)
	if err != nil {
		logger.Error("can't send start trip command to kafka:", zap.Error(err))
		return err
	}
	return nil
}

func (s *driverService) EndTrip(ctx context.Context, driverId uuid.UUID, tripId uuid.UUID) error {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	newCtx, span := tr.Start(ctx, "driver.service: EndTrip")
	defer span.End()

	ctx = zapctx.WithLogger(newCtx, logger)

	trip, err := s.r.GetTripByID(newCtx, tripId)
	if err != nil || trip == nil {
		logger.Error("can't get trip from repository", zap.Error(err))
		return err
	}
	//if trip.Status == nil || *trip.Status != domain.TripStatuses.GetStarted() {
	//	return errors.Wrap(domain.ErrAccessDenied, "trip hasn't connected with driver yet")
	//}
	if trip.DriverId != nil && *trip.DriverId != driverId.String() {
		return errors.Wrap(domain.ErrAccessDenied, "trip driver id does not match passed id")
	}
	dId := driverId.String()
	trip.DriverId = &dId
	err = s.kafkaClient.SendTripStatusUpdate(newCtx, *trip, domain.TripCommandEnd, nil)
	if err != nil {
		logger.Error("can't send end trip command to kafka:", zap.Error(err))
		return err
	}
	return nil
}

// Long poll
func (s *driverService) GetTrips(ctx context.Context, driverId uuid.UUID) ([]domain.Trip, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	newCtx, span := tr.Start(ctx, "driver.service: GetTrips")
	defer span.End()

	ctx = zapctx.WithLogger(newCtx, logger)

	trips := make([]domain.Trip, 0)
	// fot loop for 5-10 seconds with timeout:
	//
	// 		span.AddEvent("Message from Kafka")
	//		TripCreated <- Kafka
	//
	//		drivers := s.locationClient.GetDrivers(TripCreated.userLocation, radius)
	//
	//      for range drivers {
	//			if dr.id == id {
	//					trips = trips.append(TripCreated)
	//			{
	//      {
	//

	// TODO DELETE?
	//trips, err := s.locationClient.GetDrivers(newCtx, driverId)
	//if err != nil {
	//	logger.Error("driver-service: EndTrip")
	//	return nil, err
	//}

	return trips, nil
}

// UpdateTripStatus changes the status of the specified trip in the service layer
func (s *driverService) UpdateTripStatus(ctx context.Context, tripId uuid.UUID, status domain.TripStatus) error {
	log := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	newCtx, span := tr.Start(ctx, "driver.service: UpdateTripStatus")
	defer span.End()

	err := s.r.ChangeTripStatus(newCtx, tripId, status)
	if err != nil {
		log.Error("driver-service: update trip status in repository", zap.Error(err))
		return domain.FilterErrors(err)
	}

	return nil
}

// GetDrivers retrieves driver locations in the service layer
func (s *driverService) GetDrivers(ctx context.Context, driverLocation domain.LatLngLiteral, radius float32) ([]domain.DriverLocation, error) {
	log := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	newCtx, span := tr.Start(ctx, "driver.service: GetDrivers")
	defer span.End()

	// Call the repository method to get driver locations
	driverLocations, err := s.locationClient.GetDrivers(newCtx, driverLocation, radius)
	if err != nil {
		log.Error("driver-service: get drivers from repository", zap.Error(err))
		return nil, domain.FilterErrors(err)
	}

	return driverLocations, nil
}
