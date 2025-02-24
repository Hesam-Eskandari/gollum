package models

import (
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/cpp/entity"
)

type CppFederal struct {
	Year                            int     `json:"year"`
	BasicRate                       float64 `json:"basicRate"`
	FirstAdditionalRate             float64 `json:"firstAdditionalRate"`
	SecondAdditionalRate            float64 `json:"secondAdditionalRate"`
	BasicExemption                  float64 `json:"basicExemption"`
	MaxPensionableEarning           float64 `json:"maxPensionableEarning"`
	AdditionalMaxPensionableEarning float64 `json:"additionalMaxPensionableEarning"`
}

func (cpp *CppFederal) ToEntity() entity.FederalCPP {
	return entity.FederalCPP{
		Year:                            cpp.Year,
		BasicRate:                       cpp.BasicRate,
		FirstAdditionalRate:             cpp.FirstAdditionalRate,
		SecondAdditionalRate:            cpp.SecondAdditionalRate,
		BasicExemption:                  cpp.BasicExemption,
		MaxPensionableEarning:           cpp.MaxPensionableEarning,
		AdditionalMaxPensionableEarning: cpp.AdditionalMaxPensionableEarning,
	}
}
