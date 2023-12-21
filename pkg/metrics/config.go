package metrics

const (
	defaultHouseEnable  = true
	defaultHouseAddress = "DefaultAddr"
)

type Config struct {
	HouseEnable  bool   `mapstructure:"enable"`
	HouseAddress string `mapstructure:"address"`
}

func NewConfig(prefix string) *Config {
	if prefix != "" {
		prefix += "."
	}

	cfg := &Config{}
	//config.BoolVar(&cfg.HouseEnable, prefix+"house.enable", defaultHouseEnable, "description")
	//config.StringVar(&cfg.HouseAddress, prefix+"house.address", defaultHouseAddress, "description")

	return cfg
}

func NewDefaultConfig() *Config {
	return &Config{
		HouseEnable:  defaultHouseEnable,
		HouseAddress: defaultHouseAddress,
	}
}
