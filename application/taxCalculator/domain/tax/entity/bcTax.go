package entity

import (
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/brackets/entity"
)

type BcTax struct {
	BPA                 float64 // tax credit
	Brackets            []entity.TaxBracket
	CPPBasic            float64 // tax credit
	CPPFirst            float64 // tax deduction
	CPPSecond           float64 // tax deduction
	CreditsTaxReduction float64
	EIP                 float64 // tax credit
	PayableTax          float64
	TotalDeductions     float64
	TotalCredits        float64
	TotalTax            float64
	Year                int
}
