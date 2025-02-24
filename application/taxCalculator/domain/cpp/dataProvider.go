package cpp

import (
	"context"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/cpp/entity"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/errorWrap"
)

type DataProvider interface {
	GetFederalCPPAsync(ctx context.Context, year int) <-chan errorWrap.ErrorWrap[entity.FederalCPP]
}
