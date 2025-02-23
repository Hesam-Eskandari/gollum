package cpp

import (
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/cpp/entity"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/errorWrap"
)

type FederalDataProvider interface {
	GetFederalCPPAsync(year int) <-chan errorWrap.ErrorWrap[entity.FederalCPP]
}
