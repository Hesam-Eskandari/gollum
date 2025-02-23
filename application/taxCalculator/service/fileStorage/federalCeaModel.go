package fileStorage

import (
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/cea/entity"
)

type ceaFederal struct {
	Year  int     `json:"year"`
	Value float64 `json:"value"`
}

func (cea *ceaFederal) toEntity() entity.FederalCEA {
	return entity.FederalCEA{
		Amount: cea.Value,
	}
}
