package scraper

import (
	"context"
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

func (s *Scraper) StopFunc() func(context.Context) error {
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

func (s *Scraper) scrape(scrapeInterval time.Duration) {
	ctx := context.Background()
	zapctx.WithLogger(ctx, s.logger)
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(ctx, "driver.daemon.scraper: Scrape", trace.WithNewRoot())

	time.Sleep(scrapeInterval)
	s.logger.Debug("[scraper]: scrapping...")
	trips, err := s.driverService.GetTripsByStatus(ctxTrace, domain.TripStatuses.GetDriverSearch())
	if err != nil {
		s.logger.Debug("err Getting Trips By Status:", zap.Error(err))
		return
	}
	span.AddEvent("Got available trips from mongo: ", trace.WithAttributes(attribute.Int("trips_len", len(trips))))
	s.logger.Debug("[scraper] available trips from mongo:", zap.Any("trips", trips))

	defer span.End()

	for _, trip := range trips {

		drivers, err := s.driverService.GetDrivers(ctxTrace, *trip.From, domain.SearchRadius, "")
		if err != nil {
			s.logger.Debug("err Getting Drivers By Radius:", zap.Error(err))
			return
		}
		span.AddEvent("Got drivers form location:", trace.WithAttributes(attribute.Int("drivers_len", len(trips))))
		//s.logger.Debug("[scraper]  drivers in radius form location:", zap.Any("drivers", drivers))

		for i, driver := range drivers {
			if i == 0 {
				continue
			}
			str := driver.DriverId.String()
			driverCaller, ok := domain.IncomingTrips.GetTripChannel(&str)
			if ok {
				//s.logger.Debug("[scraper]: sending event", zap.String("trip_id", trip.Id.String()))
				driverCaller <- trip.Id
				span.AddEvent("Sent trip to driver:", trace.WithAttributes(attribute.String("trip_id", trip.Id.String()), attribute.String("driver_id", driver.DriverId.String())))
				//s.logger.Debug("[scraper]: event sent", zap.String("trip_id", trip.Id.String()))
			}
		}
	}
}
