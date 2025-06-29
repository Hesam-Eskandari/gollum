package main

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"time"

	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/presentation/taxCalculator/controller"
	controller2 "github.com/Hesam-Eskandari/gollum/application/taxCalculator/presentation/taxMargins/controller"
	"github.com/Hesam-Eskandari/gollum/library/environment"
	"github.com/Hesam-Eskandari/gollum/library/fileLogger"
	"github.com/Hesam-Eskandari/gollum/library/httpServer"
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
	filePath, err := buildLoggerFilePath()
	if err != nil {
		panic(err)
	}
	config := fileLogger.Config{
		PurgePeriod:     time.Hour * 24 * 7,
		CheckPeriod:     time.Minute,
		MaxFileSizeByte: 1024 * 1024 * 10,
		Filename:        filePath,
		MinLogLevel:     minLogLevel,
	}
	logger := fileLogger.New(config)
	if err = logger.Setup(); err != nil {
		panic(err)
	}
	defer func() { _ = logger.Destroy() }()
	slog.Info(fmt.Sprintf("application started in %v mode", env))
	if envErr != nil {
		slog.Warn(envErr.Error())
	}
}

func buildLoggerFilePath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	logsFilePath := path.Join(dir, ".logs/tax-calculator")
	if err = os.MkdirAll(logsFilePath, os.FileMode(0777)); err != nil {
		return "", err
	}
	return path.Join(logsFilePath, "tax-calculator.log"), nil
}
