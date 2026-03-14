package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// New returns a preconfigured logger suitable for desktop app logs.
func New() *logrus.Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	return log
}
