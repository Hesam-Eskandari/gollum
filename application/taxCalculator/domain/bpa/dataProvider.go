package bpa

import (
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/bpa/entity"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/errorWrap"
)

type FederalDataProvider interface {
	GetFederalBPAAsync(year int) <-chan errorWrap.ErrorWrap[entity.FederalBPA]
}

type BCDataProvider interface {
	GetBritishColumbiaBPAAsync(year int) <-chan errorWrap.ErrorWrap[entity.BritishColumbiaBPA]
}
