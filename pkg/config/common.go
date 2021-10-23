package config


type BaseConfig struct {
	LogLevel string
	Env      string

	KafkaConfig
	PsqlConfig
}

func NewBaseConfig() (*BaseConfig, error) {
	envLogLevel, err := GetEnv("LOG_LEVEL")
	if err != nil {
		return nil, err
	}
	envEnv, err := GetEnv("ENV")
	if err != nil {
		return nil, err
	}

	return &BaseConfig{
		LogLevel: envLogLevel,
		Env:      envEnv,
	}, nil
}
