package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todoserver"
	"todoserver/appsettings"
	storageinmem "todoserver/storage/inmem"
	transportrest "todoserver/transport/rest"

	"github.com/sirupsen/logrus"
)

func getLogger() *logrus.Logger {
	var logger = logrus.New()

	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:            true, // FIXME
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})

	// Output to stdout instead of the default stderr
	logger.SetOutput(os.Stdout)

	// log all
	logger.SetLevel(logrus.TraceLevel)

	return logger
}

func getContextLogger(l *logrus.Logger, layer, pkg string) *logrus.Entry {
	return l.WithField(layer, pkg)
}

// startHttpServer é uma goroutine para Start e Graceful Shutdown do http.Server
func startHttpServer(srv *http.Server, done chan error, logger *logrus.Logger) {
	go func() {
		err := srv.ListenAndServe()

		// sempre retornará ErrServerClosed ao chamar srv.Shutdown()
		if errors.Is(err, http.ErrServerClosed) {
			logger.Trace("HTTP Server graceful shutdown")
			return
		}

		if err != nil {
			logger.Errorf("HTTP Server Failed : %s", err)
			done <- fmt.Errorf("HTTP Server Failed : %w", err)
			return
		}

		done <- nil
	}()
}

// startSignalListener é uma goroutine para capturar sinais enviados ao processo
// e chamar o Shutdown do http.Server
func startSignalListener(srv *http.Server, done chan error, logger *logrus.Logger) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		logger.Infof("signal %s received", sig)
		logger.Info("Start HTTP Server graceful shutdown")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); errors.Is(err, context.DeadlineExceeded) {
			done <- fmt.Errorf("Timeout shutting down the server : %s", err)
			return
		} else if err != nil {
			done <- fmt.Errorf("Failure shutting down the server : %s", err)
			return
		}

		logger.Info("HTTP Server graceful shutdown OK!")
		done <- nil
	}()
}

func main() {
	logger := getLogger()

	err := run(logger)
	if err != nil {
		logger.Fatal(err)
	}
}

func run(logger *logrus.Logger) error {
	logger.Infof("Starting TO DO Service. PID: %d", os.Getpid())

	settings, err := appsettings.NewAppSettings("../settings.json")
	if err != nil {
		return fmt.Errorf("Failed to get configurations : %w", err)
	}

	//storageLogger := getContextLogger(logger, "storage", "inmem")
	//serviceLogger := getContextLogger(logger, "service", "todoservice")
	transportLogger := getContextLogger(logger, "transport", "rest")

	storage := storageinmem.NewStorageTodoInmem()
	service := todoserver.NewTodoService(storage)
	restRouter := transportrest.NewTransportRest(service, "/api/v1", transportLogger)

	logger.Infof("Listening on: %s", settings.ServerInfo.Address)

	srv := &http.Server{
		Addr:    settings.ServerInfo.Address,
		Handler: restRouter.Handler,
	}

	done := make(chan error)
	startHttpServer(srv, done, logger)
	startSignalListener(srv, done, logger)

	if err := <-done; err != nil {
		logger.Errorf("Error in run : %s", err)
		return fmt.Errorf("Error in run : %w", err)
	}

	return nil
}
