package interactor

import (
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/cpp"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/cpp/entity"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/errorWrap"
)

type FederalCPP interface {
	GetCPPBasicAsync(income float64, year int) <-chan errorWrap.ErrorWrap[float64]
	CalcCPPBasic(model entity.FederalCPP, income float64) (float64, error)
	GetCPPFirstAdditionalAsync(income float64, year int) <-chan errorWrap.ErrorWrap[float64]
	CalcCPPFirstAdditional(model entity.FederalCPP, income float64) (float64, error)
	GetCPPSecondAdditionalAsync(income float64, year int) <-chan errorWrap.ErrorWrap[float64]
	CalcCPPSecondAdditional(model entity.FederalCPP, income float64) (float64, error)
}

func NewFederalCPP(dataProvider cpp.FederalDataProvider) FederalCPP {
	return &federalCPPImpl{
		dataProvider: dataProvider,
	}
}

type federalCPPImpl struct {
	dataProvider cpp.FederalDataProvider
}

func (fed *federalCPPImpl) GetCPPBasicAsync(income float64, year int) <-chan errorWrap.ErrorWrap[float64] {
	resChan := make(chan errorWrap.ErrorWrap[float64], 1)
	go func() {
		defer close(resChan)
		res := errorWrap.ErrorWrap[float64]{}
		defer func() { resChan <- res }()
		modelWrap := <-fed.dataProvider.GetFederalCPPAsync(year)
		if modelWrap.Error != nil {
			res.Error = modelWrap.Error
			return
		}
		model := modelWrap.Value
		res.Value, res.Error = fed.CalcCPPBasic(model, income)
	}()
	return resChan
}

func (fed *federalCPPImpl) CalcCPPBasic(model entity.FederalCPP, income float64) (float64, error) {
	if err := model.Validate(); err != nil {
		return 0, err
	}
	higherValue := max(min(model.MaxPensionableEarning, income)-model.BasicExemption, 0)
	basicValue := higherValue * model.BasicRate / 100
	basicValue = float64(int(basicValue*1000)) / 1000
	return basicValue, nil
}

func (fed *federalCPPImpl) GetCPPFirstAdditionalAsync(income float64, year int) <-chan errorWrap.ErrorWrap[float64] {
	resChan := make(chan errorWrap.ErrorWrap[float64], 1)
	go func() {
		defer close(resChan)
		res := errorWrap.ErrorWrap[float64]{}
		defer func() { resChan <- res }()
		modelWrap := <-fed.dataProvider.GetFederalCPPAsync(year)
		if modelWrap.Error != nil {
			res.Error = modelWrap.Error
			return
		}
		model := modelWrap.Value
		res.Value, res.Error = fed.CalcCPPFirstAdditional(model, income)
	}()
	return resChan
}

func (fed *federalCPPImpl) CalcCPPFirstAdditional(model entity.FederalCPP, income float64) (float64, error) {
	if err := model.Validate(); err != nil {
		return 0, err
	}
	if model.Year < entity.FromYearCPPFirstAdditional {
		return 0, nil
	}
	higherValue := max(min(model.MaxPensionableEarning, income)-model.BasicExemption, 0)
	firstAdditionalValue := higherValue * model.FirstAdditionalRate / 100
	firstAdditionalValue = float64(int(firstAdditionalValue*1000)) / 1000
	return firstAdditionalValue, nil
}

func (fed *federalCPPImpl) GetCPPSecondAdditionalAsync(income float64, year int) <-chan errorWrap.ErrorWrap[float64] {
	resChan := make(chan errorWrap.ErrorWrap[float64], 1)
	go func() {
		defer close(resChan)
		res := errorWrap.ErrorWrap[float64]{}
		defer func() { resChan <- res }()
		modelWrap := <-fed.dataProvider.GetFederalCPPAsync(year)
		if modelWrap.Error != nil {
			res.Error = modelWrap.Error
			return
		}
		model := modelWrap.Value
		res.Value, res.Error = fed.CalcCPPSecondAdditional(model, income)
	}()
	return resChan
}

func (fed *federalCPPImpl) CalcCPPSecondAdditional(model entity.FederalCPP, income float64) (float64, error) {
	if err := model.Validate(); err != nil {
		return 0, err
	}
	if income <= model.MaxPensionableEarning || model.Year < entity.FromYearCPPSecondAdditional {
		return 0, nil
	}
	higherValue := max(min(model.AdditionalMaxPensionableEarning, income)-model.MaxPensionableEarning, 0)
	secondAdditionalValue := higherValue * model.SecondAdditionalRate / 100
	secondAdditionalValue = float64(int(secondAdditionalValue*1000)) / 1000
	return secondAdditionalValue, nil
}
