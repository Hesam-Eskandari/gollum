package main

import (
	"fmt"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/constants/province"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/tax/interactor"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/service/fileStorage"
)

func main() {
	dp := fileStorage.GetFileStorageInstance()
	calculator := interactor.NewFederalTax(dp)
	tax, err := calculator.Calculate(2025, province.Federal, 100000)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", tax)
}
