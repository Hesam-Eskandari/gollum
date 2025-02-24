package interactor

import (
	"context"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/bpa"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/bpa/entity"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/errorWrap"
)

type FederalBPA interface {
	GetAmountAsync(ctx context.Context, income float64, year int) <-chan errorWrap.ErrorWrap[float64]
	CalculateAmount(model entity.FederalBPA, income float64) (float64, error)
}

type federalBPAImpl struct {
	bpaAmount    float64
	dataProvider bpa.FederalDataProvider
}

func NewFederalBPA(dataProvider bpa.FederalDataProvider) FederalBPA {
	return &federalBPAImpl{
		dataProvider: dataProvider,
	}
}

func (fed *federalBPAImpl) GetAmountAsync(ctx context.Context, income float64, year int) <-chan errorWrap.ErrorWrap[float64] {
	resChan := make(chan errorWrap.ErrorWrap[float64], 1)
	go func() {
		defer close(resChan)
		res := errorWrap.ErrorWrap[float64]{}
		defer func() { resChan <- res }()
		if fed.bpaAmount != 0 {
			res.Value = fed.bpaAmount
			return
		}
		modelWrap := <-fed.dataProvider.GetFederalBPAAsync(ctx, year)
		if modelWrap.Error != nil {
			res.Error = modelWrap.Error
			return
		}
		res.Value, res.Error = fed.CalculateAmount(modelWrap.Value, income)
	}()
	return resChan
}

func (fed *federalBPAImpl) CalculateAmount(model entity.FederalBPA, income float64) (float64, error) {
	if err := model.Validate(); err != nil {
		return 0, err
	}
	if income >= model.MinBPAIncome {
		fed.bpaAmount = model.MinBPAAmount
	} else if income > model.MaxBPAIncome && income < model.MinBPAIncome {
		fed.bpaAmount = model.MaxBPAAmount -
			((income - model.MaxBPAIncome) *
				((model.MaxBPAAmount - model.MinBPAAmount) / (model.MinBPAIncome - model.MaxBPAIncome)))
	} else {
		fed.bpaAmount = model.MaxBPAAmount
	}
	fed.bpaAmount = float64(int(fed.bpaAmount*1000) / 1000)
	return fed.bpaAmount, nil
}
