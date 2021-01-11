package security

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

// nolint
type Token struct {
	AccessToken string `json:"access_token"`
	// todo: add RefreshToken
}

const Secret string = "secret" // tmp solution.

func GenerateToken(data map[string]interface{}) (*Token, error) {
	claims := jwt.MapClaims{}
	for k, v := range data {
		claims[k] = v
	}
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := unsignedToken.SignedString([]byte(Secret))
	if err != nil {
		return nil, errors.Wrap(err, "security: GenerateToken -> unsignedToken.SignedString(***)")
	}
	return &Token{AccessToken: signedToken}, nil
}

func VerifyToken(t string) (*jwt.Token, error) {
	token, err := jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("security: wrong signing method")
		}
		return []byte(Secret), nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "security: VerifyToken ->jwt.Parse(***, func)")
	}
	return token, nil
}

func ValidateToken(t string) error {
	token, err := VerifyToken(t)
	if err != nil || !token.Valid {
		return errors.New("security: unable to verify token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("security: unexpected claims type")
	}
	expiration, ok := claims["exp"]
	if !ok {
		return errors.New("security: exp is absent in claims")
	}
	val, ok := expiration.(float64)
	if !ok {
		return errors.New("security: unable to convert exp to int64")
	}
	if int64(val) < time.Now().Unix() {
		return errors.New("security: token expired")
	}
	return nil
}
