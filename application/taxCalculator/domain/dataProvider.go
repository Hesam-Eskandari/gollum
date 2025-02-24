package domain

import (
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/bpa"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/brackets"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/cea"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/cpp"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/eip"
)

type FederalDataProvider interface {
	bpa.FederalDataProvider
	cea.DataProvider
	cpp.DataProvider
	eip.DataProvider
	brackets.TaxBracketDataProvider
}

type BCDataProvider interface {
	bpa.BCDataProvider
	cpp.DataProvider
	eip.DataProvider
	brackets.TaxBracketDataProvider
}

type DataProvider interface {
	FederalDataProvider
	BCDataProvider
}
