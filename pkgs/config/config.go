package config

import (
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type LoggerSingleton struct {
	Logger log.Logger
}

type publicKey struct {
	pubKey []byte
}

var loggerConfig *LoggerSingleton
var pkey *publicKey

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

func GetPublicKey() []byte {
	if pkey == nil {
		pkey = &publicKey{
			pubKey: make([]byte, 0),
		}
	}
	return pkey.pubKey
}

func SetPublicKey(key []byte) {
	if pkey == nil {
		pkey = &publicKey{
			pubKey: key,
		}
	}
}
