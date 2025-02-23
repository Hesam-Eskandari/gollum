package entity

import "fmt"

const (
	FromYearCPPFirstAdditional  int = 2019
	FromYearCPPSecondAdditional int = 2024
)

type FederalCPP struct {
	Year                            int
	BasicRate                       float64
	FirstAdditionalRate             float64
	SecondAdditionalRate            float64
	BasicExemption                  float64
	MaxPensionableEarning           float64
	AdditionalMaxPensionableEarning float64
}

func (fed *FederalCPP) Validate() error {
	if fed.BasicRate <= 0 {
		return fmt.Errorf("cpp error: invalid basic rate \"%v\" for year \"%v\"", fed.BasicRate, fed.Year)
	}
	if fed.FirstAdditionalRate <= 0 && fed.Year >= FromYearCPPFirstAdditional {
		return fmt.Errorf("cpp error: invalid first additional rate \"%v\" for year \"%v\"", fed.FirstAdditionalRate, fed.Year)
	}
	if fed.SecondAdditionalRate <= 0 && fed.Year >= FromYearCPPSecondAdditional {
		return fmt.Errorf("cpp error: invalid second additional rate \"%v\" for year \"%v\"", fed.SecondAdditionalRate, fed.Year)
	}
	return nil
}
