package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/farhanmahmood1811/common/cerror"
)

type TokenDetail struct {
	ID        string
	AccountID string
	Role      RoleType
}

const (
	Generic  RoleType = "generic"
	Consumer RoleType = "consumer"
	Admin    RoleType = "admin"
)

var adminJwtKey []byte
var consumerJwtKey []byte
var adminCode string

// Struct that will be encoded to a JWT.
type claims struct {
	ID        string
	AccountId string
	Role      RoleType
	jwt.StandardClaims
}

func ValidateAdminCode(
	code string,
) bool {
	return code == adminCode
}

func Generate(
	id string,
	accountID string,
	role RoleType,
	expireAt time.Time,
) (string, error) {
	jwtKey := consumerJwtKey

	claims := &claims{
		ID:        id,
		AccountId: accountID,
		Role:      role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireAt.Unix(),
		},
	}
	if role == Admin {
		jwtKey = adminJwtKey
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func validate(
	token string,
	role RoleType,
) (*claims, error) {
	c := claims{}
	var err error
	if role == Admin {
		err = c.decodeAdmin(token)
	} else {
		err = c.decodeConsumer(token)
	}

	if err != nil {
		return nil, err
	}

	if c.Role != role {
		return nil, cerror.Forbidden(
			cerror.TokenInvalidCode,
			"Access Denied",
		)
	}

	if c.ExpiresAt < time.Now().Unix() {
		return nil, cerror.Unauthorised(
			cerror.TokenExpiredCode,
			"Token Expired",
		)
	}
	return &c, nil
}

func GetID(
	token string,
) (string, error) {
	c := claims{}
	err := c.decodeConsumer(token)
	if err != nil {
		return "", err
	}
	return c.Id, nil
}

func GetDecodedDetail(
	token string,
) (*TokenDetail, error) {
	c := claims{}
	err := c.decodeConsumer(token)
	if err != nil {
		return nil, err
	}
	detail := TokenDetail{
		ID:        c.ID,
		AccountID: c.AccountId,
		Role:      c.Role,
	}
	return &detail, nil
}

func (cl *claims) decodeAdmin(token string) error {
	tkn, err := jwt.ParseWithClaims(
		token,
		cl,
		func(token *jwt.Token) (interface{}, error) {
			return adminJwtKey, nil
		},
	)

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return cerror.BadRequest(
				cerror.TokenInvalidCode,
				cerror.InvalidRequestMsg,
			)
		}
		return err
	}
	if !tkn.Valid {
		return cerror.InvalidToken()
	}
	return nil
}

func (cl *claims) decodeConsumer(
	token string,
) error {
	tkn, err := jwt.ParseWithClaims(
		token,
		cl,
		func(token *jwt.Token) (interface{}, error) {
			return consumerJwtKey, nil
		},
	)

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return cerror.InvalidToken()
		}
		v, _ := err.(*jwt.ValidationError)
		if v.Errors == jwt.ValidationErrorExpired {
			return cerror.TokenExpired()
		}
		return cerror.InvalidToken()
	}
	if !tkn.Valid {
		return cerror.InvalidToken()
	}
	return nil
}
