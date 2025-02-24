package eip

import (
	"context"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/eip/entity"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/errorWrap"
)

type DataProvider interface {
	GetFederalEIPAsync(ctx context.Context, year int) <-chan errorWrap.ErrorWrap[entity.FederalEIP]
}
