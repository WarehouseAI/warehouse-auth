package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

const (
	local = "local"
	prod  = "prod"
)

func Setup(env string, log *logrus.Logger) *logrus.Logger {
	if env == prod {
		file, err := os.OpenFile("./auth.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

		if err != nil {
			log.Fatalln("❌Failed to set up the logger", err)
		}

		log.Out = file
	}

	return log
}
