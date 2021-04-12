package rest

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type authUser struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}

type appClaims struct {
	Authorized bool `json:"authorized"`
	jwt.StandardClaims
}

const (
	ACCESS_SECRET = "teste"
	EXPIRY_TIME   = time.Second * 15
)

var signingKey = []byte(ACCESS_SECRET)

func createToken(username string) (string, error) {

	claims := appClaims{
		Authorized: true,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "todo-server",
			Subject:   username,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(EXPIRY_TIME).Unix(),
			// NotBefore: endOfDay().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	var claims appClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		// O algoritmo de assinatura utilizado deve ser HMAC.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	fmt.Println(token, claims)
	return token, nil
}

// helpers

// endOfDay end of day
// func endOfDay() time.Time {
// 	now := time.Now()
// 	y, m, d := now.Date()
// 	return time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), now.Location())
// }
