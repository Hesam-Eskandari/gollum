package entity

import "fmt"

// FederalBPA is a federal tax credit reducing tax more for lower incomes and less for higher incomes
type FederalBPA struct {
	MaxBPAIncome float64 // any income bellow this value will receive MaxBPAAmount tax credit
	MinBPAIncome float64 // any income above this value will receive MinBPAAmount tax credit
	MinBPAAmount float64
	MaxBPAAmount float64
}

func (fed *FederalBPA) Validate(income float64) error {
	if income < 0 {
		return fmt.Errorf("bpaf error: invalid income: \"%v\"", income)
	}
	if fed.MaxBPAIncome <= 0 {
		return fmt.Errorf("bpaf error: invalid MaxBPAIncome: \"%v\"", fed.MaxBPAIncome)
	}
	if fed.MinBPAIncome <= 0 {
		return fmt.Errorf("bpaf error: invalid MinBPAIncome: \"%v\"", fed.MinBPAIncome)
	}
	if fed.MinBPAAmount <= 0 {
		return fmt.Errorf("bpaf error: invalid MinBPAAmount: \"%v\"", fed.MinBPAAmount)
	}
	if fed.MaxBPAAmount <= 0 {
		return fmt.Errorf("bpaf error: invalid MaxBPAAmount: \"%v\"", fed.MaxBPAAmount)
	}
	return nil
}
