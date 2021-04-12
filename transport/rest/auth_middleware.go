package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// AuthMiddleware faz a autenticação e autorização da requisição HTTP
type AuthMiddleware struct {
	logger  *logrus.Entry
	tr      *TransportRest
	authJwt *AuthJwt
}

func NewAuthMiddleware(tr *TransportRest, authJwt *AuthJwt, logger *logrus.Entry) *AuthMiddleware {
	return &AuthMiddleware{
		tr:      tr,
		logger:  logger,
		authJwt: authJwt,
	}
}

func (am *AuthMiddleware) authHandlerFunc(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// A autenticação pode ocorrer em uma ou duas etapas
		// 1 - Tenta autenticar usando o token do HEADER Authorization
		// Em caso de sucesso não tenta autenticar com Cookie
		// 2 - Em caso de erro tenta autenticar usando o Cookie de acesso

		authCookie, _ := r.Cookie("access_token")
		authHeader := r.Header.Get("Authorization")

		errAuth1 := am.authUsingHeader(authHeader)
		if errAuth1 != nil && authCookie == nil {
			am.logger.Trace("Error while authenticating using in HEADER Authentication Token. NOTE: COOKIE not provided")
			am.tr.sendErrorResponse(http.StatusUnauthorized, errAuth1.Error(), w, r)
			return
		}

		if errAuth1 != nil {
			am.logger.Trace("Error while authenticating using in HEADER Authentication Token. Trying with cookie provided")
			errAuth2 := am.authUsingCookie(*authCookie)
			if errAuth2 != nil {
				am.logger.Trace("Error while authenticating using COOKIE")
				am.tr.sendErrorResponse(http.StatusUnauthorized, fmt.Sprintf("%s , %s", errAuth1, errAuth2), w, r)
				return
			}
		}

		if errAuth1 == nil {
			am.logger.Trace("Authenticated using HEADER")
		} else {
			am.logger.Trace("Authenticated using COOKIE")
		}

		h(w, r, ps)
	}
}

func (am *AuthMiddleware) authUsingCookie(cookie http.Cookie) error {
	return am.auth(cookie.Value)
}

func (am *AuthMiddleware) authUsingHeader(headerValue string) error {
	authHeaderValue, err := parseBearerToken(headerValue)
	if err != nil {
		return err
	}

	return am.auth(authHeaderValue)
}

func (am *AuthMiddleware) auth(tokenString string) error {
	if tokenString == "" {
		return fmt.Errorf("required Bearer Authorization")
	}
	_, err := am.authJwt.verifyToken(tokenString)
	if err != nil {
		return err
	}

	return nil
}

// Helpers

// parseBearerToken parses an HTTP Bearer Token
func parseBearerToken(auth string) (string, error) {
	const prefix = "Bearer "
	// Case insensitive prefix match.
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return "", fmt.Errorf("invalid Bearer Authentication")
	}
	return auth[len(prefix):], nil
}
