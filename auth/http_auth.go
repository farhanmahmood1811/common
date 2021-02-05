package auth

import (
	"context"
	"net/http"
	"strconv"

	"github.com/farhanmahmood1811/common/cerror"
	"github.com/farhanmahmood1811/common/logs"
)

type Entity struct{}

type RoleType string
type ProfileID string
type AccountID string

func NewClient(
	consumerKey, adminKey string,
) *Entity {
	if consumerKey == "" || adminKey == "" {
		panic("Auth Token Keys Missing")
	}
	adminCode = adminKey
	adminJwtKey = []byte(adminKey)
	consumerJwtKey = []byte(consumerKey)
	return &Entity{}
}

type authContextKey string

const (
	authIdContextKey authContextKey = "authContextKey"
)

func (ultronAuth Entity) ProcessAuthentication(
	next http.Handler,
	role RoleType,
) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Header["Authorization"] == nil {
			err := cerror.Unauthorised(
				cerror.TokenInvalidCode,
				"Token Missing",
			)
			cerror.SendHttpError(
				request.Context(),
				responseWriter,
				err,
			)
			return
		}

		bearerToken := request.Header["Authorization"][0]
		if bearerToken == "" {
			err := cerror.Unauthorised(
				cerror.TokenInvalidCode,
				"Token Missing",
			)
			cerror.SendHttpError(
				request.Context(),
				responseWriter,
				err,
			)
			return
		}

		if request.Header["X-Client-Time"] == nil {
			logs.GetClient().Error("Request Time Missing")
			err := cerror.BadRequest(
				cerror.InvalidClientTimeCode,
				"Request Time Missing",
			)
			cerror.SendHttpError(
				request.Context(),
				responseWriter,
				err,
			)
			return
		}

		xClientTime := request.Header["X-Client-Time"][0]
		if xClientTime == "" {
			logs.GetClient().Error("Request Time Missing")
			err := cerror.BadRequest(
				cerror.InvalidClientTimeCode,
				"Request Time Missing",
			)
			cerror.SendHttpError(
				request.Context(),
				responseWriter,
				err,
			)
			return
		}

		_, err := strconv.Atoi(xClientTime)
		if err != nil {
			logs.GetClient().Error("Invalid Client Time")
			err := cerror.BadRequest(
				cerror.InvalidClientTimeCode,
				"Invalid Client Time",
			)
			cerror.SendHttpError(
				request.Context(),
				responseWriter,
				err,
			)
			return
		}

		if request.Header["X-Device-Id"] == nil {
			err := cerror.BadRequest(
				cerror.MissingDeviceIdCode,
				"Device CounterID Missing",
			)
			cerror.SendHttpError(
				request.Context(),
				responseWriter,
				err,
			)
			return
		}

		claims, err := validate(
			bearerToken,
			role,
		)
		if err != nil {
			cerror.SendHttpError(
				request.Context(),
				responseWriter,
				err,
			)
			return
		}
		tokenDetail := TokenDetail{
			ID:        claims.ID,
			AccountID: claims.AccountId,
			Role:      claims.Role,
		}
		ctx := context.WithValue(
			request.Context(),
			authIdContextKey,
			tokenDetail,
		)
		request = request.WithContext(ctx)

		next.ServeHTTP(responseWriter, request)
	}
}

func IDFromCtx(
	ctx context.Context,
) string {
	tokenDetail, _ := ctx.Value(authIdContextKey).(TokenDetail)
	return tokenDetail.ID
}

func GetAccountIDFromCtx(
	ctx context.Context,
) string {
	tokenDetail, _ := ctx.Value(authIdContextKey).(TokenDetail)
	return tokenDetail.AccountID
}

func profileIdFromCtx(
	ctx context.Context,
) string {
	tokenDetail, _ := ctx.Value(authIdContextKey).(TokenDetail)
	return tokenDetail.ID
}

func accountIDFromCtx(
	ctx context.Context,
) string {
	tokenDetail, _ := ctx.Value(authIdContextKey).(TokenDetail)
	return tokenDetail.AccountID
}

type UserReference struct {
	ProfileID ProfileID
	AccountID AccountID
}

func GetUserReference(
	requestContext context.Context,
) UserReference {
	return UserReference{
		ProfileID: ProfileID(profileIdFromCtx(requestContext)),
		AccountID: AccountID(accountIDFromCtx(requestContext)),
	}
}