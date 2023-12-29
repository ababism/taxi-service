package kafkaproducer

type Config struct {
	Broker string `mapstructure:"broker"`
	Topic  string `mapstructure:"topic"`
}
