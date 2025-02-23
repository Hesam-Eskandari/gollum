package entity

import "fmt"

// BritishColumbiaBPA is a tax credit for bc residents that reduces tax for all incomes
type BritishColumbiaBPA struct {
	Amount float64
}

func (bc *BritishColumbiaBPA) Validate() error {
	if bc.Amount <= 0 {
		return fmt.Errorf("bpa bc error: invalid basic personal amount: \"%v\"", bc.Amount)
	}
	return nil
}
