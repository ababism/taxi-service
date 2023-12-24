package tripsscraper

import "time"

type Config struct {
	ScrapeInterval string `mapstructure:"scrape_interval"`
}

func (c Config) GetScrapeInterval() (time.Duration, error) {
	return time.ParseDuration(c.ScrapeInterval)
}
