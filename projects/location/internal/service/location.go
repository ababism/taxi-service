package service

import (
	"context"
	"github.com/google/uuid"
	"gitlab/ArtemFed/mts-final-taxi/projects/location/internal/domain"
	"gitlab/ArtemFed/mts-final-taxi/projects/location/internal/service/adapters"
	global "go.opentelemetry.io/otel"
)

var _ adapters.LocationService = &locationService{}

type locationService struct {
	r adapters.LocationRepository
}

func NewLocationService(ur adapters.LocationRepository) adapters.LocationService {
	return &locationService{r: ur}
}

func (s *locationService) GetDrivers(ctx context.Context, lat float32, lng float32, radius float32) ([]domain.Driver, error) {
	tr := global.Tracer(domain.TracerName)
	ctxTrace, span := tr.Start(ctx, "[LocationService-GetDrivers]")
	defer span.End()

	err := validateCoordinates(lat, lng)
	if err != nil {
		return nil, err
	}

	drivers, err := s.r.GetDrivers(ctxTrace, lat, lng, radius)
	if err != nil {
		return nil, err
	}

	return drivers, err
}

func (s *locationService) UpdateDriverLocation(ctx context.Context, driverId uuid.UUID, lat float32, lng float32) error {
	tr := global.Tracer(domain.TracerName)
	ctxTrace, span := tr.Start(ctx, "[LocationService-UpdateDriverLocation]")
	defer span.End()

	err := validateCoordinates(lat, lng)
	if err != nil {
		return err
	}

	err = s.r.UpdateDriverLocation(ctxTrace, driverId, lat, lng)
	if err != nil {
		return err
	}

	return err
}

func validateCoordinates(lat float32, lng float32) error {
	if lat < -90 || lat > 90 {
		return domain.ErrBadLatitude
	}

	if lng < -180 || lng > 180 {
		return domain.ErrBadLongitude
	}

	return nil
}
