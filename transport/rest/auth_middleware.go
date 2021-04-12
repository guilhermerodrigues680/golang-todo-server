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
	next   http.Handler
	logger *logrus.Entry
}

func NewAuthMiddleware(next http.Handler, logger *logrus.Entry) *AuthMiddleware {
	return &AuthMiddleware{next: next, logger: logger}
}

func (m *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.logger.Trace("Auth...")

	if !strings.HasSuffix(r.URL.Path, "/authenticate") {
		err := auth(r.Header.Get("Authorization"))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, err)
			return
		}

		// w.WriteHeader(http.StatusForbidden)
	}

	m.next.ServeHTTP(w, r)
}

func authHandlerFunc(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// A autenticação ocorre em duas etapas
		// 1 - Tenta autentica usar o HEADER Authorization
		// 2 - Em caso de erro tenta autenticar usando o Cookie de acesso

		authCookie, _ := r.Cookie("access_token")
		authHeader := r.Header.Get("Authorization")

		errAuth1 := authUsingHeader(authHeader)
		if errAuth1 != nil && authCookie == nil {
			logrus.Info("Erro ao autenticar usando HEADER, COOKIE não fornecido")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Erro: %s", errAuth1)
			return
		}

		if errAuth1 != nil {
			logrus.Info("Erro ao autenticar usando HEADER, tentando com COOKIE fornecido")
			errAuth2 := authUsingCookie(*authCookie)
			if errAuth2 != nil {
				logrus.Info("Erro ao autenticar usando HEADER E COOKIE")
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, "Erros: %s , %s", errAuth1, errAuth2)
				return
			}
		}

		if errAuth1 == nil {
			logrus.Info("Autenticando usando HEADER")
		} else {
			logrus.Info("Autenticando usando COOKIE")
		}

		h(w, r, ps)
	}
}

// Helper

// parseBearerToken parses an HTTP Bearer Token
func parseBearerToken(auth string) (string, error) {
	const prefix = "Bearer "
	// Case insensitive prefix match.
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return "", fmt.Errorf("invalid Bearer Authentication")
	}
	return auth[len(prefix):], nil
}

func auth(tokenString string) error {
	if tokenString == "" {
		return fmt.Errorf("required Bearer Authorization")
	}
	_, err := verifyToken(tokenString)
	if err != nil {
		return err
	}

	return nil
}

func authUsingCookie(cookie http.Cookie) error {
	return auth(cookie.Value)
}

func authUsingHeader(headerValue string) error {
	authHeaderValue, err := parseBearerToken(headerValue)
	if err != nil {
		return err
	}

	return auth(authHeaderValue)
}
