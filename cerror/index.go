package cerror

import (
	"context"
	"encoding/json"
	"net/http"

	"go.elastic.co/apm"
)

type Entity struct {
	HttpStatusCode int
	Message        string
	Code           Code
}

func (ue *Entity) Error() string {
	return ue.Message
}

func ValidationError(msg string) error {
	return &Entity{
		HttpStatusCode: http.StatusBadRequest,
		Message:        msg,
		Code:           ParamMissingCode,
	}
}

func InternalServerError() error {
	return &Entity{
		HttpStatusCode: http.StatusInternalServerError,
		Message:        "Something Went Wrong",
		Code:           ServerErrorCode,
	}
}

func Unauthorised(code Code, message string) error {
	return &Entity{
		HttpStatusCode: http.StatusUnauthorized,
		Message:        message,
		Code:           code,
	}
}

func Forbidden(code Code, message string) error {
	return &Entity{
		HttpStatusCode: http.StatusForbidden,
		Message:        message,
		Code:           code,
	}
}

func BadRequest(code Code, message string) error {
	return &Entity{
		HttpStatusCode: http.StatusBadRequest,
		Message:        message,
		Code:           code,
	}
}

func InvalidToken() error {
	return &Entity{
		HttpStatusCode: http.StatusBadRequest,
		Message:        "Auth Token Invalid",
		Code:           TokenInvalidCode,
	}
}

func TokenExpired() error {
	return &Entity{
		HttpStatusCode: http.StatusBadRequest,
		Message:        "Auth Token Expired",
		Code:           TokenExpiredCode,
	}
}

func SendHttpError(ctx context.Context, w http.ResponseWriter, err error) {
	apm.CaptureError(ctx, err).Send()
	httpStatusCode := http.StatusInternalServerError
	code := ServerErrorCode
	message := err.Error()

	if errVal, ok := err.(*Entity); ok {
		httpStatusCode = errVal.HttpStatusCode
		message = errVal.Message
		code = errVal.Code
	}

	resp := map[string]interface{}{
		"message": message,
		"code":    code,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	json.NewEncoder(w).Encode(resp)
	return
}

