package kafkaproducer

type Config struct {
	Brokers []string `mapstructure:"brokers"`
	Topic   string   `mapstructure:"topic"`
	IdGroup string   `mapstructure:"id_group"`
}
