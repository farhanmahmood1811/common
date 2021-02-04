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

func MuxObservability(router *mux.Router) {
	apmgorilla.Instrument(router)
}

func MongoObservability(client *options.ClientOptions) {
	client.SetMonitor(apmmongo.CommandMonitor())
}

func LogObservability(client *logrus.Logger) {
	apm.DefaultTracer.SetLogger(client)
	client.AddHook(&apmlogrus.Hook{})
}
