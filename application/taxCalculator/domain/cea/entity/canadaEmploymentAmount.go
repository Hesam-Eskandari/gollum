package entity

import "fmt"

type FederalCEA struct {
	Amount float64
}

func (fed *FederalCEA) Validate() error {
	if fed.Amount <= 0 {
		return fmt.Errorf("cea error: invalid cea: \"%v\"", fed.Amount)
	}
	return nil
}
