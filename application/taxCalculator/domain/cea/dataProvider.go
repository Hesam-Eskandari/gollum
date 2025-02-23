package cea

import (
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/cea/entity"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/errorWrap"
)

type DataProvider interface {
	GetFederalCEAAsync(year int) <-chan errorWrap.ErrorWrap[entity.FederalCEA]
}
