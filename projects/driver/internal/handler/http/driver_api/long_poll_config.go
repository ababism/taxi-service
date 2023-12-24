package driver_api

import "time"

type Config struct {
	LongPollTimeout string `mapstructure:"long_poll_timeout"`
}

func (c Config) GetLongPollTimeout() (time.Duration, error) {
	return time.ParseDuration(c.LongPollTimeout)
}
