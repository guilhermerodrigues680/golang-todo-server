package rest

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

type appClaims struct {
	Authorized bool `json:"authorized"`
	jwt.StandardClaims
}

type AuthJwt struct {
	accessSecret string
	expiryTime   time.Duration
	logger       *logrus.Entry
}

func NewAuthJwt(logger *logrus.Entry) *AuthJwt {
	const (
		ACCESS_SECRET = "teste"
		EXPIRY_TIME   = time.Second * 15
	)

	return &AuthJwt{
		accessSecret: ACCESS_SECRET,
		expiryTime:   EXPIRY_TIME,
		logger:       logger,
	}
}

func (ajwt *AuthJwt) createToken(username string) (string, *appClaims, error) {

	expiresAt := time.Now().Add(ajwt.expiryTime)
	claims := appClaims{
		Authorized: true,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "todo-server",
			Subject:   username,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expiresAt.Unix(),
			// NotBefore: endOfDay().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(ajwt.accessSecret))
	if err != nil {
		return "", nil, err
	}

	ajwt.logger.Infof("iat: %s, eat: %s", time.Now(), expiresAt)

	return tokenString, &claims, nil
}

func (ajwt *AuthJwt) verifyToken(tokenString string) (*jwt.Token, error) {
	var claims appClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		// O algoritmo de assinatura utilizado deve ser HMAC.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(ajwt.accessSecret), nil
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
