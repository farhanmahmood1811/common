package middleware

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/farhanmahmood1811/common/auth"
)

type ContextKey string

// ContextKey key name for context value
const (
	ReqIDContextKey ContextKey = "requestId"
)

var ultronAuthEntity *auth.Entity

// InitializeMiddleware sets chttp client for middleware
func Use(
	ultronAuth *auth.Entity,
) {
	ultronAuthEntity = ultronAuth
}

func AuthWrapMiddleware(next http.HandlerFunc, roles auth.RoleType) http.HandlerFunc {
	next = recoverHandler(next)
	next = reqIDMiddleware(next)
	next = addLoggerMiddleware(next)
	next = ultronAuthEntity.ProcessAuthentication(next, roles)

	return next
}

// UnAuthWrapMiddleware wraps and applies multiple middleware without auth
func UnAuthWrapMiddleware(next http.HandlerFunc) http.HandlerFunc {
	next = recoverHandler(next)
	next = reqIDMiddleware(next)
	next = addLoggerMiddleware(next)
	return next
}

// recoverHandler middleware to recover from
func recoverHandler(
	next http.HandlerFunc,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				//logger := logs.GetFromContext(r.Context())
				//logger.Client.Info(err)
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				fmt.Printf("Stack trace: %s", buf)
				serverErr := ultronerror.InternalServerError()
				ultronerror.SendHttpError(r.Context(), w, serverErr)
			}
		}()
		next.ServeHTTP(w, r)
	}
}

// contextLoggerMiddleware add context to req logger
func reqIDMiddleware(
	next http.HandlerFunc,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profileId := auth.IDFromCtx(r.Context())
		if profileId == "" {
			profileId = "generic"
		} else {
			contextData := appendValueToLogger(r.Context(), "profile_id", profileId)
			r.WithContext(contextData)
		}

		reqID := generateRequestID(profileId)

		contextData := appendValueToLogger(r.Context(), "request_id", reqID)
		r.WithContext(contextData)

		contextData = context.WithValue(contextData, ReqIDContextKey, reqID)
		r = r.WithContext(contextData)

		next.ServeHTTP(w, r)
	}
}

func addLoggerMiddleware(
	next http.HandlerFunc,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httpLogger := logs.HttpLogger{}
		httpLogger.Populate(r)
		contextData := logs.GetUpdateLoggerContext(r.Context(), httpLogger)

		r = r.WithContext(contextData)
		next.ServeHTTP(w, r)
	}
}

// generateRequestId create request id
func generateRequestID(
	profileId string,
) string {
	epochTime := time.Now().UnixNano() / 1000000
	reqID := "ultron::uid:" + profileId + "::t:" + strconv.FormatInt(epochTime, 10)
	return reqID
}

func appendValueToLogger(
	ctx context.Context,
	key, value string,
) context.Context {
	logger := logs.GetFromContext(ctx)
	logger.UpdateContextField(key, value)
	contextData := logs.GetUpdateLoggerContext(ctx, logger)
	return contextData
}
