package main

import (
	"fmt"
	"net/http"
	"os"
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

func main() {
	logger := getLogger()

	err := run(logger)
	if err != nil {
		logger.Fatal(err)
	}
}

func run(logger *logrus.Logger) error {
	logger.Info("Starting TO DO Service")

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

	err = http.ListenAndServe(settings.ServerInfo.Address, restRouter.Handler)
	if err != nil {
		return fmt.Errorf("HTTP Server Failed : %w", err)
	}

	return nil
}
