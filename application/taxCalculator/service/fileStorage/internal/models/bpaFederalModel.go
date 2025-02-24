package models

import "github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/bpa/entity"

type BpaFederal struct {
	Year         int     `json:"year"`
	MinBpaAmount float64 `json:"minBpaAmount"`
	MaxBpaAmount float64 `json:"maxBpaAmount"`
	MinBpaIncome float64 `json:"minBpaIncome"`
	MaxBpaIncome float64 `json:"maxBpaIncome"`
}

func (bp *BpaFederal) ToEntity() entity.FederalBPA {
	return entity.FederalBPA{
		MaxBPAAmount: bp.MaxBpaAmount,
		MinBPAAmount: bp.MinBpaAmount,
		MaxBPAIncome: bp.MaxBpaIncome,
		MinBPAIncome: bp.MinBpaIncome,
	}
}
