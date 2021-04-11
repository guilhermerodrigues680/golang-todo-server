package rest

import (
	"net"
	"net/http"

	"github.com/sirupsen/logrus"
)

// LoggingMiddleware faz o log da requisição HTTP
type LoggingMiddleware struct {
	next   http.Handler
	logger *logrus.Entry
}

func NewLoggingMiddleware(next http.Handler, logger *logrus.Entry) *LoggingMiddleware {
	return &LoggingMiddleware{next: next, logger: logger}
}

func (m *LoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	remoteIp, _, _ := net.SplitHostPort(r.RemoteAddr)
	forwardedFor := r.Header.Get("X-Forwarded-For")
	m.logger.Infof("Start: %s - %s '%s %s'", forwardedFor, remoteIp, r.Method, r.RequestURI)
	lw := NewLoggingResponseWriter(w)
	m.next.ServeHTTP(lw, r)
	m.logger.Infof("End:   %s - %s '%s %s' %d", forwardedFor, remoteIp, r.Method, r.RequestURI, lw.statusCode)
}

// loggingResponseWriter faz o log do código HTTP da resposta
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	// WriteHeader(int) não é chamado se nossa resposta retornar implicitamente 200 OK, então
	// configura-se por default este status code.
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lw *loggingResponseWriter) WriteHeader(code int) {
	lw.statusCode = code
	lw.ResponseWriter.WriteHeader(code)
}
