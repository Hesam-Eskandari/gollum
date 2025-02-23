package entity

import (
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/brackets/entity"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/constants/province"
)

type FederalTax struct {
	BPA                 float64 // tax credit
	CEA                 float64 // tax credit
	CPPBasic            float64 // tax credit
	CPPFirst            float64 // tax deduction
	CPPSecond           float64 // tax deduction
	EIP                 float64 // tax credit
	TotalDeductions     float64
	TotalCredits        float64
	CreditsTaxReduction float64
	Brackets            []entity.TaxBracket
	Year                int
	Province            province.Province
	TotalTax            float64
	PayableTax          float64
}
