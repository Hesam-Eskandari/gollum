package interactor

import (
	"context"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/cea"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/errorWrap"
)

type FederalCEA interface {
	GetAmountAsync(ctx context.Context, year int) <-chan errorWrap.ErrorWrap[float64]
}

func NewFederalCEA(dataProvider cea.DataProvider) FederalCEA {
	return &federalCEAImpl{
		dataProvider: dataProvider,
	}
}

type federalCEAImpl struct {
	amount       float64
	dataProvider cea.DataProvider
}

func (fed *federalCEAImpl) GetAmountAsync(ctx context.Context, year int) <-chan errorWrap.ErrorWrap[float64] {
	resChan := make(chan errorWrap.ErrorWrap[float64], 1)
	go func() {
		defer close(resChan)
		res := errorWrap.ErrorWrap[float64]{}
		defer func() { resChan <- res }()
		if fed.amount != 0 {
			res.Value = fed.amount
			return
		}
		modelWrap := <-fed.dataProvider.GetFederalCEAAsync(ctx, year)
		if modelWrap.Error != nil {
			res.Error = modelWrap.Error
			return
		}
		if err := modelWrap.Value.Validate(); err != nil {
			res.Error = err
			return
		}
		model := modelWrap.Value
		fed.amount = model.Amount
		res.Value = fed.amount
	}()
	return resChan
}
