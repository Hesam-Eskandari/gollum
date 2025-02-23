package eip

import (
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/eip/entity"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/errorWrap"
)

type DataProvider interface {
	GetFederalEIPAsync(year int) <-chan errorWrap.ErrorWrap[entity.FederalEIP]
}
