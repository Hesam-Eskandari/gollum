package interactor

import (
	"context"
	"iter"
	"math"
	"slices"

	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/brackets"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/brackets/entity"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/constants/province"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/errorWrap"
)

type TaxBracketCalculator interface {
	GetCalculatedTaxAsync(ctx context.Context, taxableIncome float64, year int, pr province.Province) <-chan errorWrap.ErrorWrap[float64]
	ApplyTaxBrackets(taxableIncome float64, brackets []entity.TaxBracket) float64
	GetAllMarginalBrackets(ctx context.Context, year int, pr province.Province) (iter.Seq2[int, entity.TaxBracket], error)
}

func NewTaxBracketCalculator(dataProvider brackets.TaxBracketDataProvider) TaxBracketCalculator {
	return &taxBracketCalculatorImpl{
		dataProvider: dataProvider,
	}
}

type taxBracketCalculatorImpl struct {
	dataProvider brackets.TaxBracketDataProvider
}

func (tb *taxBracketCalculatorImpl) GetAllMarginalBrackets(ctx context.Context, year int, pr province.Province) (iter.Seq2[int, entity.TaxBracket], error) {
	fedChan := tb.dataProvider.GetTaxBracketsAsync(ctx, year, province.Federal)
	prChan := tb.dataProvider.GetTaxBracketsAsync(ctx, year, pr)
	fedWrap := <-fedChan
	if fedWrap.Error != nil {
		return nil, fedWrap.Error
	}
	prWrap := <-prChan
	if prWrap.Error != nil {
		return nil, prWrap.Error
	}
	return func(yield func(int, entity.TaxBracket) bool) {
		var fedBracket, prBracket entity.TaxBracket
		var low, high, rate float64
		var i int
		for f, p := 0, 0; f < len(fedWrap.Value) || p < len(prWrap.Value); {
			low = high
			if f >= len(fedWrap.Value) {
				high = prWrap.Value[p].High
				rate = prWrap.Value[p].Rate + fedBracket.Rate
				p++
			} else if p >= len(prWrap.Value) {
				high = fedWrap.Value[f].High
				rate = fedWrap.Value[f].Rate + prBracket.Rate
				f++
			} else {
				fedBracket = fedWrap.Value[f]
				prBracket = prWrap.Value[p]
				rate = fedBracket.Rate + prBracket.Rate
				if fedBracket.High < prBracket.High {
					high = fedBracket.High
					f++
				} else {
					high = prBracket.High
					p++
				}
			}
			if !yield(i, entity.TaxBracket{
				Low:  low,
				High: high,
				Rate: rate,
			}) {
				break
			}
			if high == math.Inf(1) {
				break
			}
			i++
		}
	}, nil
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
