package cerror

const (
	InvalidRequestMsg   string = "Invalid Request"
	UnauthorisedRequest string = "Unauthorised Request"
)

type Code string

const (
	InvalidRequestCode    Code = "CL_BAD_REQUEST"
	ParamMissingCode      Code = "CL_PARAM_MISSING"
	InvalidIDCode         Code = "CL_INVALID_ID"
	InvalidClientTimeCode Code = "CL_INVALID_CLIENT_TIME"
	MissingDeviceIdCode   Code = "CL_MISSING_DEVICE_ID"
	TokenExpiredCode      Code = "AUTH_TOKEN_EXPIRED"
	TokenInvalidCode      Code = "AUTH_TOKEN_INVALID"
	NotAuthorisedCode     Code = "AUTH_NO_PERMISSION"
	ServerErrorCode       Code = "SERVER_ERROR"
)
