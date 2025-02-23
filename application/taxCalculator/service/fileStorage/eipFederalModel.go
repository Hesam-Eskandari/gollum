package fileStorage

import (
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/eip/entity"
)

type eipFederal struct {
	Year                              int     `json:"year"`
	Rate                              float64 `json:"rate"`
	MaxInsurableEarning               float64 `json:"maxInsurableEarning"`
	EmployerEmployeeContributionRatio float64 `json:"employerEmployeeContributionRatio"`
}

func (eip *eipFederal) toEntity() entity.FederalEIP {
	return entity.FederalEIP{
		Rate:                              eip.Rate,
		MaxInsurableEarning:               eip.MaxInsurableEarning,
		EmployerEmployeeContributionRatio: eip.EmployerEmployeeContributionRatio,
	}
}
