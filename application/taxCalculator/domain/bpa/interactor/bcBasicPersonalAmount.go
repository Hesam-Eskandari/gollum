package interactor

import (
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/bpa"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/errorWrap"
)

type BritishColumbiaBPA interface {
	GetAmountAsync(year int) <-chan errorWrap.ErrorWrap[float64]
}

func NewBritishColumbiaBPA(dataProvider bpa.BCDataProvider) BritishColumbiaBPA {
	return &bcBpaImpl{
		dataProvider: dataProvider,
	}
}

type bcBpaImpl struct {
	amount       float64
	dataProvider bpa.BCDataProvider
}

func (bpa *bcBpaImpl) GetAmountAsync(year int) <-chan errorWrap.ErrorWrap[float64] {
	resChan := make(chan errorWrap.ErrorWrap[float64], 1)
	go func() {
		defer close(resChan)
		res := errorWrap.ErrorWrap[float64]{}
		defer func() { resChan <- res }()
		if bpa.amount != 0 {
			res.Value = bpa.amount
			return
		}
		modelWrap := <-bpa.dataProvider.GetBritishColumbiaBPAAsync(year)
		if modelWrap.Error != nil {
			res.Error = modelWrap.Error
			return
		}
		if err := modelWrap.Value.Validate(); err != nil {
			res.Error = err
			return
		}
		res.Value = modelWrap.Value.Amount
	}()
	return resChan
}
