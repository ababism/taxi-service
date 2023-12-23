package mylogger

import (
	"github.com/TheZeroSlave/zapsentry"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

const (
	envDebug      = "dev"
	envProduction = "prod"
)

func InitLogger(cfg *Config, appName string) (*zap.Logger, error) {
	var levelZap, levelSentry zapcore.Level
	err := levelZap.UnmarshalText([]byte(cfg.LevelLogger))
	if err != nil {
		log.Printf("Zap logs level with value=%s not initiolized", levelZap)
		return nil, err
	}

	var cfgZap zap.Config
	if cfg.Env == envDebug {
		cfgZap = zap.NewDevelopmentConfig()
	} else {
		cfgZap = zap.NewProductionConfig()
	}
	cfgZap.Level = zap.NewAtomicLevelAt(levelZap)
	cfgZap.OutputPaths = cfg.OutputPaths
	cfgZap.ErrorOutputPaths = cfg.ErrorOutputPaths

	logger, err := cfgZap.Build()
	if err != nil {
		log.Println("Zap logger build failed")
		return nil, err
	}

	err = levelSentry.UnmarshalText([]byte(cfg.LevelSentry))
	if err != nil {
		log.Printf("Sentry logs level with value=%s not initiolized", levelZap)
		return nil, err
	}
	cfgSentry := zapsentry.Configuration{
		Level: levelSentry,
		Tags: map[string]string{
			"environment": cfg.Env,
			"app":         appName,
		},
	}
	core, err := zapsentry.NewCore(
		cfgSentry,
		zapsentry.NewSentryClientFromDSN(cfg.SentryDSN),
	)
	if err != nil {
		log.Println("Zapsentry NewCore not initiolized")
		return nil, err
	}

	logger = zapsentry.AttachCoreToLogger(core, logger)

	return logger, nil
}
