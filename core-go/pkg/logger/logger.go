package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func New(env string) *logrus.Logger {
	l := logrus.New()
	l.Out = os.Stdout
	l.SetFormatter(&logrus.JSONFormatter{})
	if env == "production" {
		l.SetLevel(logrus.InfoLevel)
	} else {
		l.SetLevel(logrus.DebugLevel)
	}
	return l
}
