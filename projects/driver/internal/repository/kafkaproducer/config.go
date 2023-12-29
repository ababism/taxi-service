package kafkaproducer

type Config struct {
	Broker  string `mapstructure:"brokers"`
	Topic   string `mapstructure:"topic"`
	IdGroup string `mapstructure:"id_group"`
}
