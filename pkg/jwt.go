package pkg

import (
	"fmt"
	"log"
	"os"
	"time"

	msgerr "github.com/MGT06/e-wallet_backend/internal/message"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	ID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func NewJWTClaims(user_id int64, expired int) *JWTClaims {
	return &JWTClaims{
		ID:        user_id,
		Issuer:    os.Getenv("JWT_ISSUER"),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expired) * time.Minute)),
	}
}

func (j *JWTClaims) GeneretToken() (string, error) {
	jwtKey := os.Getenv("JWT_SECREET_KEY")
	if jwtKey == "" {
		log.Printf("[PkgJwt.GeneretToken] Getenv JWT_SECREET_KEY : %v", msgerr.ErrMissingKey)
		return "", msgerr.ErrMissingKey
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, j)
	return token.SignedString([]byte(jwtKey))
}

func (j *JWTClaims) DecodeToken(tokenString string) error {
	jwtKey := os.Getenv("JWT_SECREET_KEY")
	if jwtKey == "" {
		log.Printf("[PkgJwt.DecodeToken] Getenv JWT_SECREET_KEY : %v", msgerr.ErrMissingKey)
		return msgerr.ErrMissingKey
	}

	jwtToken, err := jwt.ParseWithClaims(tokenString, j, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("[PkgJwt.DecodeToken] ParseWithClaims : %v", msgerr.ErrUnexpectedSign)
			return nil, fmt.Errorf("%w: %v", msgerr.ErrUnexpectedSign, t.Header["alg"])
		}
		return []byte(jwtKey), nil
	})
	if err != nil {
		return err
	}

	if !jwtToken.Valid {
		return jwt.ErrTokenInvalidClaims
	}

	iss, err := jwtToken.Claims.GetIssuer()
	if err != nil {
		log.Printf("[PkgJwt.DecodeToken] Claims.GetIssuer : %v", err)
		return err
	}
	if iss != os.Getenv("JWT_ISSUER") {
		log.Printf("[PkgJwt.GeneretToken] Getenv JWT_ISSUER : %v", jwt.ErrTokenInvalidIssuer)
		return jwt.ErrTokenInvalidIssuer
	}

	return nil
}
