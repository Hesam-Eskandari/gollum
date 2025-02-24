package models

import "github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/bpa/entity"

type BpaBc struct {
	Year  int     `json:"year"`
	Value float64 `json:"value"`
}

func (bpa *BpaBc) ToEntity() entity.BritishColumbiaBPA {
	return entity.BritishColumbiaBPA{
		Amount: bpa.Value,
	}
}
