package interactor

import (
	"context"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/eip"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/eip/entity"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/errorWrap"
)

type FederalEIP interface {
	GetEmployeeContributionAsync(ctx context.Context, income float64, year int) <-chan errorWrap.ErrorWrap[float64]
	CalcEmployeeContribution(model entity.FederalEIP, income float64) (float64, error)
	GetEmployerContributionAsync(ctx context.Context, income float64, year int) <-chan errorWrap.ErrorWrap[float64]
	CalcEmployerContribution(model entity.FederalEIP, income float64) (float64, error)
}

func NewFederalEIP(dataProvider eip.DataProvider) FederalEIP {
	return &federalEIPImpl{
		dataProvider: dataProvider,
	}
}

type federalEIPImpl struct {
	dataProvider eip.DataProvider
}

func (fed *federalEIPImpl) GetEmployeeContributionAsync(ctx context.Context, income float64, year int) <-chan errorWrap.ErrorWrap[float64] {
	resChan := make(chan errorWrap.ErrorWrap[float64], 1)
	go func() {
		defer close(resChan)
		res := errorWrap.ErrorWrap[float64]{}
		defer func() { resChan <- res }()
		modelChan := <-fed.dataProvider.GetFederalEIPAsync(ctx, year)
		if modelChan.Error != nil {
			res.Error = modelChan.Error
			return
		}
		model := modelChan.Value
		res.Value, res.Error = fed.CalcEmployeeContribution(model, income)
	}()
	return resChan
}

func (fed *federalEIPImpl) CalcEmployeeContribution(model entity.FederalEIP, income float64) (float64, error) {
	if err := model.Validate(); err != nil {
		return 0, nil
	}
	lowerValue := min(income, model.MaxInsurableEarning)
	employeeContribution := lowerValue * model.Rate / 100
	employeeContribution = float64(int(employeeContribution*1000)) / 1000
	return employeeContribution, nil
}

func (fed *federalEIPImpl) GetEmployerContributionAsync(ctx context.Context, income float64, year int) <-chan errorWrap.ErrorWrap[float64] {
	resChan := make(chan errorWrap.ErrorWrap[float64], 1)
	go func() {
		defer close(resChan)
		res := errorWrap.ErrorWrap[float64]{}
		defer func() { resChan <- res }()
		modelWrap := <-fed.dataProvider.GetFederalEIPAsync(ctx, year)
		if modelWrap.Error != nil {
			res.Error = modelWrap.Error
			return
		}
		if err := modelWrap.Value.Validate(); err != nil {
			res.Error = err
			return
		}
		res.Value, res.Error = fed.CalcEmployerContribution(modelWrap.Value, income)
	}()
	return resChan
}

func (fed *federalEIPImpl) CalcEmployerContribution(model entity.FederalEIP, income float64) (float64, error) {
	employeeContribution, err := fed.CalcEmployeeContribution(model, income)
	if err != nil {
		return 0, err
	}
	employerContribution := employeeContribution * model.EmployerEmployeeContributionRatio
	employerContribution = float64(int(employerContribution*1000)) / 1000
	return employeeContribution, nil
}
