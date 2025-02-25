package main

import (
	"context"
	"fmt"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/tax/interactor"
	controller2 "github.com/Hesam-Eskandari/gollum/application/taxCalculator/presentation/taxBrackets/controller"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/presentation/taxCalculator/controller"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/service/fileStorage"
	"github.com/Hesam-Eskandari/gollum/internal/httpServer"
)

func main() {
	server := httpServer.NewServer(":51483", false)
	_ = server.AddController(controller.NewTaxCalculatorController())
	_ = server.AddController(controller2.NewTaxBracketsController())
	if err := <-server.Launch(); err != nil {
		panic(err)
	}
}

func runDemo() {
	ctx := context.Background()
	dp := fileStorage.GetFileStorageInstance()
	federalCalculator := interactor.NewFederalTax(dp)
	fedTax, err := federalCalculator.Calculate(ctx, 2025, 100000)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", fedTax)

	bcCalculator := interactor.NewBcTax(dp)
	bcTax, err := bcCalculator.Calculate(ctx, 2025, 100000)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", bcTax)
}
