package cea

import (
	"context"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/cea/entity"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/errorWrap"
)

type DataProvider interface {
	GetFederalCEAAsync(ctx context.Context, year int) <-chan errorWrap.ErrorWrap[entity.FederalCEA]
}
