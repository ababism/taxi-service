package scraper

import (
	"context"
	"github.com/google/uuid"
	"github.com/juju/zaputil/zapctx"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/service/adapters"
	global "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"time"
)

type Scraper struct {
	stop          chan bool
	driverService adapters.DriverService
	logger        *zap.Logger
}

func NewScraper(logger *zap.Logger, driverService adapters.DriverService) *Scraper {
	return &Scraper{
		logger:        logger,
		driverService: driverService,
		stop:          make(chan bool)}
}

func (s *Scraper) stopCallback(ctx context.Context) error {
	s.stop <- true
	return nil
}

func (s Scraper) StopFunc() func(context.Context) error {
	return s.stopCallback
}

func (s *Scraper) Start(scrapeInterval time.Duration) {
	go func() {
		stop := s.stop
		go func() {
			for {
				s.scrape(scrapeInterval)
			}
		}()
		<-stop
	}()
}

func generateRequestID() string {
	id := uuid.New()
	return id.String()
}

func WithRequestID(ctx context.Context) context.Context {
	requestID := generateRequestID()
	return context.WithValue(ctx, domain.KeyRequestID, requestID)
}

func (s *Scraper) scrape(scrapeInterval time.Duration) {
	ctx := context.Background()

	requestIdCtx := WithRequestID(ctx)
	ctxLogger := zapctx.WithLogger(requestIdCtx, s.logger)

	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(ctxLogger, "driver.daemon.scraper: Scrape", trace.WithNewRoot())

	time.Sleep(scrapeInterval)

	s.logger.Debug("[scraper]: scrapping...")
	trips, err := s.driverService.GetTripsByStatus(ctxTrace, domain.TripStatuses.GetDriverSearch())
	if err != nil {
		s.logger.Debug("err Getting trips By Status:", zap.Error(err))
		return
	}
	span.AddEvent("Got available trips from mongo: ", trace.WithAttributes(attribute.Int("trips_len", len(trips))))

	defer span.End()

	for _, trip := range trips {
		drivers, err := s.driverService.GetDrivers(ctxTrace, *trip.From, domain.SearchRadius, "")
		if err != nil {
			s.logger.Debug("err Getting Drivers By Radius:", zap.Error(err))
			return
		}
		span.AddEvent("Got drivers form location for trip:", trace.WithAttributes(attribute.String("trip_id", trip.Id.String()), attribute.Int("drivers_len", len(drivers))))
		s.logger.Debug("[driver.scraper]  got drivers in radius from location:", zap.String("trip_id", trip.Id.String()), zap.Any("drivers", drivers))

		for _, driver := range drivers {
			drID := driver.DriverId.String()

			_ = domain.AvailableTripEvents.SendTrip(drID, trip.Id)
		}
	}
}
