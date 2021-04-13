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
	"todoapp"
	"todoapp/appsettings"
	storagepostgres "todoapp/storage/postgres"
	transportrest "todoapp/transport/rest"

	"github.com/sirupsen/logrus"
)

// getLogger retorna uma instância do logger com configurações pré-definidas
func getLogger(isProduction bool) *logrus.Logger {
	var logger = logrus.New()

	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:            !isProduction,
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})

	if !isProduction {
		// log tudo
		logger.SetLevel(logrus.TraceLevel)
		// Output to stdout instead of the default stderr
		logger.SetOutput(os.Stdout)
		return logger
	}

	// log somente da severidade info ou acima
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})
	file, err := os.OpenFile("todoapp.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.SetOutput(file)
		//multiWriter := io.MultiWriter(os.Stdout, file)
		//logger.SetOutput(multiWriter)
	} else {
		logger.Info("Failed to log to file, using Stdout")
		logger.SetOutput(os.Stdout)
	}

	return logger
}

// getContextLogger retorna um logger com contexto
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

		timeout := 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
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
	settings, err := appsettings.NewAppSettings("../settings.json")
	if err != nil {
		logrus.Fatalf("Failed to get configurations : %s", err)
	}

	logger := getLogger(settings.Environment.IsProduction)

	logger.Infof("Starting TO-DO App, PID: %d", os.Getpid())
	logger.Infof("App Mode: '%s', Is Production: %t", settings.Environment.AppMode, settings.Environment.IsProduction)

	err = run(settings, logger)
	if err != nil {
		logger.Fatal(err)
	}
}

// run é responsável por inicializar e finalizar a aplicação
func run(settings *appsettings.AppSettings, logger *logrus.Logger) error {
	// storageInmemLogger := getContextLogger(logger, "storage", "inmem")
	storagePgsqlLogger := getContextLogger(logger, "storage", "postgres")
	serviceLogger := getContextLogger(logger, "service", "todoservice")
	transportLogger := getContextLogger(logger, "transport", "rest")

	storagePgsqlService, err := storagepostgres.NewPostgresService(
		settings.StorageCredentials.DBHost,
		settings.StorageCredentials.DBPort,
		settings.StorageCredentials.DBName,
		settings.StorageCredentials.DBUser,
		settings.StorageCredentials.DBPassword,
		storagePgsqlLogger)

	if err != nil {
		return err
	}

	storagePgsqlService.MigrateTables()
	storageTodo := storagepostgres.NewPostgresTodo(storagePgsqlService, storagePgsqlLogger)

	// storageInmem := storageinmem.NewStorageTodoInmem(storageInmemLogger)
	// service := todoapp.NewTodoService(storageInmem, serviceLogger)
	service := todoapp.NewTodoService(storageTodo, serviceLogger)

	restRouter, err := transportrest.NewTransportRest(service, transportLogger)
	if err != nil {
		return err
	}

	logger.Infof("Listening on: %s", settings.ServerInfo.Address)

	srv := &http.Server{
		Addr:    settings.ServerInfo.Address,
		Handler: restRouter,
	}

	done := make(chan error)
	startHttpServer(srv, done, logger)
	startSignalListener(srv, done, logger)

	if err := <-done; err != nil {
		errWrapped := fmt.Errorf("Error in run : %w", err)
		logger.Error(errWrapped)
		return errWrapped
	}

	return nil
}
