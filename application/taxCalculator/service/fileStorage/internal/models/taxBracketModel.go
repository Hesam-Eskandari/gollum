package models

import (
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/brackets/entity"
	"math"
	"slices"
)

type Bracket struct {
	Year int     `json:"year"`
	Rate float64 `json:"rate"`
	Low  float64 `json:"low"`
}

func MapBracketsToEntity(brk []Bracket) []entity.TaxBracket {
	slices.SortFunc(brk, func(a, b Bracket) int {
		return int(a.Low - b.Low)
	})
	entities := make([]entity.TaxBracket, 0, len(brk))
	for i, b := range brk {
		high := math.MaxFloat64
		if i < len(brk)-1 {
			high = brk[i+1].Low
		}
		entities = append(entities, entity.TaxBracket{
			High: high,
			Rate: b.Rate,
			Low:  b.Low,
		})
	}
	return entities
}
