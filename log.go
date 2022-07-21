package main

import (
	"github.com/cdfmlr/crud/log"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func InitLogger() {
	// crud packages will use this logger automatically
	logger = log.NewLogger(
		log.WithLevel(log.Level(GlobalConfig.LogLevel)),
		log.WithReportCaller(false),
		log.WithHook(log.RequestIDHook()),
	)

	logger.Info("logger is ready")
}
