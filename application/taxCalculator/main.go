package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/presentation/taxCalculator/controller"
	controller2 "github.com/Hesam-Eskandari/gollum/application/taxCalculator/presentation/taxMargins/controller"

	"github.com/Hesam-Eskandari/gollum/internal/environment"
	"github.com/Hesam-Eskandari/gollum/internal/fileLogger"
	"github.com/Hesam-Eskandari/gollum/internal/httpServer"
)

func main() {
	setup()
	runServers()
}

func runServers() {
	server := httpServer.NewServer(":51483", false)
	_ = server.AddController(controller.NewTaxCalculatorController())
	_ = server.AddController(controller2.NewTaxMarginsController())
	if err := <-server.Launch(); err != nil {
		panic(err)
	}
}

func setup() {
	envErr := environment.AutoSetEnvironment()
	env := environment.GetEnvironment()
	minLogLevel := slog.LevelInfo
	if env == environment.Develop {
		minLogLevel = slog.LevelDebug
	}
	config := fileLogger.Config{
		PurgePeriod:     time.Hour * 24 * 7,
		CheckPeriod:     time.Minute,
		MaxFileSizeByte: 1024 * 1024 * 10,
		Filename:        "./tax-calculator.log",
		MinLogLevel:     minLogLevel,
	}
	logger := fileLogger.New(config)
	if err := logger.Setup(); err != nil {
		panic(err)
	}
	defer func() { _ = logger.Destroy() }()
	slog.Info(fmt.Sprintf("application started in %v mode", env))
	if envErr != nil {
		slog.Warn(envErr.Error())
	}
}
