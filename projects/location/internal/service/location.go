package service

import (
	"context"
	"github.com/juju/zaputil/zapctx"
	openapitypes "github.com/oapi-codegen/runtime/types"
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
	tr := global.Tracer("gitlab/ArtemFed/mts-final-taxi")
	ctxTrace, span := tr.Start(ctx, "location-service: GetDrivers Service")
	defer span.End()

	logger := zapctx.Logger(ctx)

	drivers, err := s.r.GetDrivers(ctxTrace, lat, lng, radius)
	if err != nil {
		logger.Error("location-service: get drivers in radius from repository")
		return nil, domain.FilterErrors(err)
	}

	return drivers, err
}

func (s *locationService) UpdateDriverLocation(ctx context.Context, driverId openapitypes.UUID, lat float32, lng float32) error {
	tr := global.Tracer("gitlab/ArtemFed/mts-final-taxi")
	ctxTrace, span := tr.Start(ctx, "location-service: UpdateDriverLocation Service")
	defer span.End()

	logger := zapctx.Logger(ctx)

	err := s.r.UpdateDriverLocation(ctxTrace, driverId, lat, lng)
	if err != nil {
		logger.Error("location-service: update driver location in radius from repository")
		return domain.FilterErrors(err)
	}

	return err
}
