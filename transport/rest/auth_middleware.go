package rest

import (
	"fmt"
	"net/http"
	"strings"

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
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := parseBearerToken(auth)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "%s", err)
			return
		}

		_, err = verifyToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "%s", err)
			return
		}

		// w.WriteHeader(http.StatusForbidden)
	}

	m.next.ServeHTTP(w, r)
}

// parseBearerToken parses an HTTP Bearer Token
func parseBearerToken(auth string) (string, error) {
	const prefix = "Bearer "
	// Case insensitive prefix match.
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return "", fmt.Errorf("invalid Bearer Authentication")
	}
	return auth[len(prefix):], nil
}
