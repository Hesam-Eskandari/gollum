package interactor

import (
	"context"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/brackets"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/brackets/entity"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/constants/province"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/errorWrap"
	"slices"
)

type TaxBracketCalculator interface {
	GetCalculatedTaxAsync(ctx context.Context, taxableIncome float64, year int, pr province.Province) <-chan errorWrap.ErrorWrap[float64]
	ApplyTaxBrackets(taxableIncome float64, brackets []entity.TaxBracket) float64
}

func NewTaxBracketCalculator(dataProvider brackets.TaxBracketDataProvider) TaxBracketCalculator {
	return &taxBracketCalculatorImpl{
		dataProvider: dataProvider,
	}
}

type taxBracketCalculatorImpl struct {
	dataProvider brackets.TaxBracketDataProvider
}

func (tb *taxBracketCalculatorImpl) GetCalculatedTaxAsync(ctx context.Context, taxableIncome float64, year int, pr province.Province) <-chan errorWrap.ErrorWrap[float64] {
	resChan := make(chan errorWrap.ErrorWrap[float64], 1)
	go func() {
		defer close(resChan)
		res := errorWrap.ErrorWrap[float64]{}
		defer func() { resChan <- res }()
		modelEWrap := <-tb.dataProvider.GetTaxBracketsAsync(ctx, year, pr)
		if modelEWrap.Error != nil {
			res.Error = modelEWrap.Error
			return
		}
		res.Value = tb.ApplyTaxBrackets(taxableIncome, modelEWrap.Value)
	}()
	return resChan
}

func (tb *taxBracketCalculatorImpl) ApplyTaxBrackets(taxableIncome float64, brackets []entity.TaxBracket) float64 {
	slices.SortFunc(brackets, func(a, b entity.TaxBracket) int { return int(a.Low - b.Low) })
	value := 0.0
	for _, bracket := range brackets {
		if taxableIncome < bracket.Low {
			break
		}
		lowerValue := min(bracket.High, taxableIncome)
		value += bracket.Rate * (lowerValue - bracket.Low) / 100.0
	}
	return value
}
