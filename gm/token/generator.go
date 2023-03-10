package token

import (
	"github.com/dgrijalva/jwt-go"
)

var Bearer = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJvdmVybG9yZCJ9.N2LpIoR0fS4rUCF7ul0NL75b5sy-NCNc_2sI7ZJ9_yU"

// GenerateBearerToken 生成固定 bearer token ，所有请求只验证此token
func GenerateBearerToken(SecretKey []byte, issuer string) (string, error) {
	claims := jwt.StandardClaims{
		Issuer: issuer,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}
