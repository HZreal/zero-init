package jwts

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"time"
)

type JwtPayload struct {
	UserId      int64  `json:"userId"`
	Username    string `ctx:"username"`
	RoleId      int64  `json:"roleId"`   // 角色 ID
	RoleType    int64  `json:"roleType"` // 角色类型
	PrivilegeId int64  `json:"privilegeId"`
	Lang        string `json:"lang"`
}

type CustomClaims struct {
	JwtPayload
	jwt.RegisteredClaims
}

func GenJwt(user JwtPayload, accessSecret string, expires int64) (string, error) {
	claim := CustomClaims{
		JwtPayload: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(expires))),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(accessSecret))
}

func ParseToken(tokens, accessSecret string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokens, &CustomClaims{}, Secret(accessSecret))
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("that's not even a token")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("token is expired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("token not active yet")
			} else {
				return nil, errors.New("couldn't handle this token")
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("couldn't handle this token")
}

func Secret(accessSecret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(accessSecret), nil
	}
}
