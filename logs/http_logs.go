package logs

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmlogrus"
)

type ContextKey string

const (
	loggerContextKey ContextKey = "httpLogger"
)

type HttpLogger struct {
	Client *logrus.Entry
	Fields logrus.Fields
}

func (h *HttpLogger) Populate(
	r *http.Request,
) {
	logFields := logrus.Fields{}
	logFields["http_method"] = r.Method
	logFields["remote_addr"] = r.RemoteAddr
	logFields["uri"] = r.RequestURI
	log := GetClient().WithFields(
		apmlogrus.TraceContext(
			r.Context(),
		),
	)
	h.Client = log.WithFields(logFields)
	h.Fields = logFields
}
