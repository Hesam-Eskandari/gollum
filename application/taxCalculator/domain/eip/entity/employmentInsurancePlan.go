package entity

import "fmt"

type FederalEIP struct {
	MaxInsurableEarning               float64
	Rate                              float64
	EmployerEmployeeContributionRatio float64
	eiEmployeeContribution            float64
	eiEmployerContribution            float64
}

func (eip *FederalEIP) Validate() error {
	if eip.MaxInsurableEarning <= 0 {
		return fmt.Errorf("ei error: invalid max insurrable earcning: \"%v\"", eip.MaxInsurableEarning)
	}
	if eip.Rate < 0 || eip.Rate > 100 {
		return fmt.Errorf("ei error: invalid rate: \"%v\"", eip.Rate)
	}
	if eip.EmployerEmployeeContributionRatio < 0 {
		return fmt.Errorf("ei error: invalid employer contribution ratio: \"%v\"", eip.EmployerEmployeeContributionRatio)
	}
	return nil
}
