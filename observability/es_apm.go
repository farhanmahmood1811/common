package observability

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmgorilla"
	"go.elastic.co/apm/module/apmlogrus"
	"go.elastic.co/apm/module/apmmongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func muxObservability(router *mux.Router) {
	apmgorilla.Instrument(router)
}

func mongoObservability(client *options.ClientOptions) {
	client.SetMonitor(apmmongo.CommandMonitor())
}

func logObservability(client *logrus.Logger) {
	apm.DefaultTracer.SetLogger(client)
	client.AddHook(&apmlogrus.Hook{})
}

func InitialiseEsApm(
	router *mux.Router,
	mongoClientOptions *options.ClientOptions,
	loggerClient *logrus.Logger,
) {
	if router != nil {
		muxObservability(router)
	}
	if mongoClientOptions != nil {
		mongoObservability(mongoClientOptions)
	}
	if loggerClient != nil {
		logObservability(loggerClient)
	}
}

