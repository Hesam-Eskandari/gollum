package bpa

import (
	"context"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/bpa/entity"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/errorWrap"
)

type FederalDataProvider interface {
	GetFederalBPAAsync(ctx context.Context, year int) <-chan errorWrap.ErrorWrap[entity.FederalBPA]
}

type BCDataProvider interface {
	GetBritishColumbiaBPAAsync(ctx context.Context, year int) <-chan errorWrap.ErrorWrap[entity.BritishColumbiaBPA]
}
