package models

import (
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/cea/entity"
)

type CeaFederal struct {
	Year  int     `json:"year"`
	Value float64 `json:"value"`
}

func (cea *CeaFederal) ToEntity() entity.FederalCEA {
	return entity.FederalCEA{
		Amount: cea.Value,
	}
}
