package logs

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/56-Secure/common/observability"
)

// Client variable
var client *logrus.Logger

// NewClient Client
func NewClient(observabilityEnabled bool) {
	client = &logrus.Logger{
		Out:   os.Stderr,
		Hooks: make(logrus.LevelHooks),
		Level: logrus.InfoLevel,
		Formatter: &logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "@timestamp",
				logrus.FieldKeyLevel: "log.level",
				logrus.FieldKeyMsg:   "message",
				logrus.FieldKeyFunc:  "function.name", // non-ECS
			},
		},
	}

	client.SetReportCaller(true)

	if observabilityEnabled == true {
		observability.LogObservability(client)
	}

}

func GetClient() *logrus.Logger {
	return client
}

