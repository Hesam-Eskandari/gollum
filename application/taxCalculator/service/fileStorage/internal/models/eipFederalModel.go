package models

import (
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/eip/entity"
)

type EipFederal struct {
	Year                              int     `json:"year"`
	Rate                              float64 `json:"rate"`
	MaxInsurableEarning               float64 `json:"maxInsurableEarning"`
	EmployerEmployeeContributionRatio float64 `json:"employerEmployeeContributionRatio"`
}

func (eip *EipFederal) ToEntity() entity.FederalEIP {
	return entity.FederalEIP{
		Rate:                              eip.Rate,
		MaxInsurableEarning:               eip.MaxInsurableEarning,
		EmployerEmployeeContributionRatio: eip.EmployerEmployeeContributionRatio,
	}
}
