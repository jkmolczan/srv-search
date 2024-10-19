package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"

	"github.com/jkmolczan/srv-search/pkg/numbers"
	pkgHttp "github.com/jkmolczan/srv-search/pkg/numbers/adapter/http"
	"github.com/jkmolczan/srv-search/pkg/numbers/infra/storage"
)

func main() {
	e := echo.New()
	logger := e.Logger

	configPath, err := filepath.Abs("config.yaml")
	if err != nil {
		logger.Fatalf("failed to get executable path for config file: %v", err)
	}

	inputDataPath, err := filepath.Abs("input.txt")
	if err != nil {
		logger.Fatalf("failed to get executable path for input data file: %v", err)
	}

	// Load configuration
	config, err := loadConfig(configPath)
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	setupLogging(config.LogLevel, logger)

	numbersStorage, err := storage.NewNumbersStorage(inputDataPath)
	if err != nil {
		logger.Fatalf("failed to load data to numbers storage: %v", err)
	}

	searchService := numbers.NewSearchService(numbersStorage)
	searchHandler := pkgHttp.NewSearchHandler(searchService, logger)

	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.File("/docs/api/swagger.yaml", fmt.Sprintf("%s/swagger.yaml", "./docs/api"))

	e.HTTPErrorHandler = pkgHttp.ErrorHandler
	pkgHttp.SetSearchNumbersRoutes(e, searchHandler)

	err = e.Start(fmt.Sprintf(":%d", config.Port))
	if err != nil {
		logger.Fatalf("server failed: %v", err)
	}
}

func setupLogging(level string, logger echo.Logger) {
	level = strings.ToLower(level)
	switch level {
	case "debug":
		logger.SetLevel(log.DEBUG)
	case "info":
		logger.SetLevel(log.INFO)
	case "error":
		logger.SetLevel(log.ERROR)
	default:
		logger.SetLevel(log.INFO)
	}
}
