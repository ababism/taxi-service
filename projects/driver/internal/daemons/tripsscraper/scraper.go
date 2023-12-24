package tripsscraper

import (
	"context"
	"github.com/juju/zaputil/zapctx"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/service/adapters"
	"go.uber.org/zap"
	"time"
)

type Scraper struct {
	stop          chan bool
	driverService adapters.DriverService
	log           zap.Logger
}

func NewScraper(logger zap.Logger, driverService adapters.DriverService) *Scraper {
	return &Scraper{
		log:           logger,
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
		for {
			stop := s.stop
			go s.scrape(scrapeInterval)
			<-stop
		}
	}()
}

func (s *Scraper) scrape(scrapeInterval time.Duration) {

	time.Sleep(scrapeInterval)

	ctx := context.Background()
	zapctx.WithLogger(ctx, &s.log)
	trips, err := s.driverService.GetTripsByStatus(ctx, domain.TripStatuses.GetDriverSearch())
	if err != nil {
		s.log.Debug("err Getting Trips By Status:", zap.Error(err))
		return
	}
	for _, trip := range trips {
		drivers, err := s.driverService.GetDrivers(ctx, *trip.From, domain.SearchRadius)
		if err != nil {
			s.log.Debug("err Getting Drivers By Radius:", zap.Error(err))
			return
		}
		for _, driver := range drivers {
			str := driver.DriverId.String()
			driverCaller, ok := domain.IncomingTrips.GetTripChannel(&str)
			if ok {
				driverCaller <- trip.Id
			}
		}
	}
}
