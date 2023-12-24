package kafkaconsumer

type Config struct {
	Brokers []string `mapstructure:"brokers"`
	Topic   string   `mapstructure:"topic"`
	IdGroup string   `mapstructure:"id_group"`
	// MinBytes Minimal Batch size that kafka will read
	MinBytes int `mapstructure:"min_bytes"`
	MaxBytes int `mapstructure:"max_bytes"`
}
