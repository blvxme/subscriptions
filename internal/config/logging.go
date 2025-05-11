package config

import (
	"github.com/sirupsen/logrus"
	"os"
)

func ConfigureLogger(logger *logrus.Logger) {
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})
}
