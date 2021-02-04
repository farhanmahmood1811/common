package response

import (
	"encoding/json"
	"net/http"
)

func Send(
	responseWriter http.ResponseWriter,
	statusCode int,
	response interface{},
) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)
	json.NewEncoder(responseWriter).Encode(response)
	return
}

