package config

import (
	"os"

	"github.com/go-kit/kit/log/level"
)

func LogAndTerminateOnError(err error, action string) {
	if err != nil {
		level.Error(GetLogger().Logger).Log("action", action, "err", err)
		os.Exit(1)
	}
}
