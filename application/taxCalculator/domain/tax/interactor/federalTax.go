package interactor

import (
	"context"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain"
	interactor3 "github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/bpa/interactor"
	interactor5 "github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/brackets/interactor"
	interactor4 "github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/cea/interactor"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/constants/province"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/cpp/interactor"
	interactor2 "github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/eip/interactor"
	entity5 "github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/tax/entity"
)

type FederalTax interface {
	Calculate(ctx context.Context, year int, income float64) (entity5.FederalTax, error)
}

func NewFederalTax(dataProvider domain.FederalDataProvider) FederalTax {
	return &federalTaxImpl{
		dataProvider:         dataProvider,
		cppCalculator:        interactor.NewFederalCPP(dataProvider),
		eipCalculator:        interactor2.NewFederalEIP(dataProvider),
		bpaCalculator:        interactor3.NewFederalBPA(dataProvider),
		ceaCalculator:        interactor4.NewFederalCEA(dataProvider),
		taxBracketCalculator: interactor5.NewTaxBracketCalculator(dataProvider),
	}
}

type federalTaxImpl struct {
	cppCalculator        interactor.FederalCPP
	eipCalculator        interactor2.FederalEIP
	bpaCalculator        interactor3.FederalBPA
	ceaCalculator        interactor4.FederalCEA
	taxBracketCalculator interactor5.TaxBracketCalculator
	dataProvider         domain.FederalDataProvider
}

func (fed *federalTaxImpl) Calculate(ctx context.Context, year int, income float64) (entity5.FederalTax, error) {
	federalTax, err := fed.getData(ctx, year, income)
	if err != nil {
		return federalTax, err
	}
	federalTax.TotalDeductions = float64(int((federalTax.CPPFirst+federalTax.CPPSecond)*1000)) / 1000
	taxableIncome := float64(int((income-federalTax.TotalDeductions)*1000)) / 1000
	federalTax.TotalTax = float64(int(fed.taxBracketCalculator.ApplyTaxBrackets(taxableIncome, federalTax.Brackets))*1000) / 1000
	federalTax.TotalCredits = float64(int((federalTax.CEA+federalTax.EIP+federalTax.BPA+federalTax.CPPBasic)*1000)) / 1000
	//TODO: check if only the first bracket should be applied for tax credit calculation
	federalTax.CreditsTaxReduction = float64(int(fed.taxBracketCalculator.ApplyTaxBrackets(federalTax.TotalCredits, federalTax.Brackets)*1000)) / 1000
	federalTax.PayableTax = float64(int((federalTax.TotalTax-federalTax.CreditsTaxReduction)*1000)) / 1000

	return federalTax, nil
}

func (fed *federalTaxImpl) getData(ctx context.Context, year int, income float64) (entity5.FederalTax, error) {
	bpaEWrapChan := fed.bpaCalculator.GetAmountAsync(ctx, income, year)
	ceaEWrapChan := fed.ceaCalculator.GetAmountAsync(ctx, year)
	cppBasicEWrapChan := fed.cppCalculator.GetCPPBasicAsync(ctx, income, year)
	cppFirstEWrapChan := fed.cppCalculator.GetCPPFirstAdditionalAsync(ctx, income, year)
	cppSecondEWrapChan := fed.cppCalculator.GetCPPSecondAdditionalAsync(ctx, income, year)
	eipEWrapChan := fed.eipCalculator.GetEmployeeContributionAsync(ctx, income, year)
	bracketsEWrapChan := fed.dataProvider.GetTaxBracketsAsync(ctx, year, province.Federal)
	federalTax := entity5.FederalTax{Year: year}
	bpaEWrap := <-bpaEWrapChan
	if bpaEWrap.Error != nil {
		return federalTax, bpaEWrap.Error
	}
	federalTax.BPA = bpaEWrap.Value
	ceaEWrap := <-ceaEWrapChan
	if ceaEWrap.Error != nil {
		return federalTax, ceaEWrap.Error
	}
	federalTax.CEA = ceaEWrap.Value
	cppBasicEWrap := <-cppBasicEWrapChan
	if cppBasicEWrap.Error != nil {
		return federalTax, cppBasicEWrap.Error
	}
	federalTax.CPPBasic = cppBasicEWrap.Value
	cppFirstEWrap := <-cppFirstEWrapChan
	if cppFirstEWrap.Error != nil {
		return federalTax, cppFirstEWrap.Error
	}
	federalTax.CPPFirst = cppFirstEWrap.Value
	cppSecondEWrap := <-cppSecondEWrapChan
	if cppSecondEWrap.Error != nil {
		return federalTax, cppSecondEWrap.Error
	}
	federalTax.CPPSecond = cppSecondEWrap.Value
	eipEWrap := <-eipEWrapChan
	if eipEWrap.Error != nil {
		return federalTax, eipEWrap.Error
	}
	federalTax.EIP = eipEWrap.Value
	bracketsEWrap := <-bracketsEWrapChan
	if bracketsEWrap.Error != nil {
		return federalTax, bracketsEWrap.Error
	}
	federalTax.Brackets = bracketsEWrap.Value
	return federalTax, nil
}
