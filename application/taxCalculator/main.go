package main

import (
	"fmt"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/tax/interactor"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/service/fileStorage"
)

func main() {
	dp := fileStorage.GetFileStorageInstance()
	federalCalculator := interactor.NewFederalTax(dp)
	fedTax, err := federalCalculator.Calculate(2025, 100000)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", fedTax)

	bcCalculator := interactor.NewBcTax(dp)
	bcTax, err := bcCalculator.Calculate(2025, 100000)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", bcTax)
}
