package ginit

import (
	"github.com/sid-sun/arche-api/config"
	"go.uber.org/zap"
)

func Logger(cfg *config.Config) (*zap.Logger, error) {
	var err error
	var logger *zap.Logger

	if cfg.GetEnv() == "dev" {
		logger, err = zap.NewDevelopmentConfig().Build()
	} else {
		logger, err = zap.NewProductionConfig().Build()
	}

	if err != nil {
		return nil, err
	}
	return logger, nil
}
