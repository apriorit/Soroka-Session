package config

import (
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type LoggerSingleton struct {
	Logger log.Logger
}

var loggerConfig *LoggerSingleton

func GetLogger() *LoggerSingleton {
	if loggerConfig == nil {
		var logger log.Logger
		{
			logger = log.NewLogfmtLogger(os.Stdout)
			logger = level.NewFilter(logger, level.AllowAll())
			logger = log.With(logger, "timestamp", log.DefaultTimestampUTC)
			logger = log.With(logger, "caller", log.DefaultCaller)
			logger = log.With(logger, "service", "sessions")
		}
		loggerConfig = &LoggerSingleton{
			Logger: logger,
		}
	}
	return loggerConfig
}
