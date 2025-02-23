package domain

import (
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/bpa"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/brackets"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/cea"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/cpp"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/eip"
)

type FederalDataProvider interface {
	bpa.DataProvider
	cea.DataProvider
	cpp.FederalDataProvider
	eip.DataProvider
	brackets.TaxBracketDataProvider
}

type DataProvider interface {
	FederalDataProvider
}
