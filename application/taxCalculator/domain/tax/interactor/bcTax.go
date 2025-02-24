package interactor

import (
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/bpa/interactor"
	interactor2 "github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/brackets/interactor"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/constants/province"
	interactor4 "github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/cpp/interactor"
	interactor3 "github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/eip/interactor"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/tax/entity"
)

type BCTax interface {
	Calculate(year int, income float64) (entity.BcTax, error)
}

func NewBcTax(dataProvider domain.BCDataProvider) BCTax {
	return &bcTaxImpl{
		dataProvider:         dataProvider,
		bpaCalculator:        interactor.NewBritishColumbiaBPA(dataProvider),
		taxBracketCalculator: interactor2.NewTaxBracketCalculator(dataProvider),
		eipCalculator:        interactor3.NewFederalEIP(dataProvider),
		cppCalculator:        interactor4.NewFederalCPP(dataProvider),
	}
}

type bcTaxImpl struct {
	dataProvider         domain.BCDataProvider
	bpaCalculator        interactor.BritishColumbiaBPA
	eipCalculator        interactor3.FederalEIP
	cppCalculator        interactor4.FederalCPP
	taxBracketCalculator interactor2.TaxBracketCalculator
}

func (bc *bcTaxImpl) Calculate(year int, income float64) (entity.BcTax, error) {
	bcTax, err := bc.getData(year, income)
	if err != nil {
		return bcTax, err
	}
	bcTax.TotalDeductions = float64(int((bcTax.CPPFirst+bcTax.CPPSecond)*1000)) / 1000
	taxableIncome := float64(int((income-bcTax.TotalDeductions)*1000)) / 1000
	bcTax.TotalTax = float64(int(bc.taxBracketCalculator.ApplyTaxBrackets(taxableIncome, bcTax.Brackets))*1000) / 1000
	bcTax.TotalCredits = float64(int((bcTax.EIP+bcTax.BPA+bcTax.CPPBasic)*1000)) / 1000
	//TODO: check if only the first bracket should be applied for tax credit calculation
	bcTax.CreditsTaxReduction = float64(int(bc.taxBracketCalculator.ApplyTaxBrackets(bcTax.TotalCredits, bcTax.Brackets)*1000)) / 1000
	bcTax.PayableTax = float64(int((bcTax.TotalTax-bcTax.CreditsTaxReduction)*1000)) / 1000
	return bcTax, nil
}

func (bc *bcTaxImpl) getData(year int, income float64) (entity.BcTax, error) {
	bpaEWrapChan := bc.bpaCalculator.GetAmountAsync(year)
	cppBasicEWrapChan := bc.cppCalculator.GetCPPBasicAsync(income, year)
	cppFirstEWrapChan := bc.cppCalculator.GetCPPFirstAdditionalAsync(income, year)
	cppSecondEWrapChan := bc.cppCalculator.GetCPPSecondAdditionalAsync(income, year)
	eipEWrapChan := bc.eipCalculator.GetEmployeeContributionAsync(income, year)
	bracketsEWrapChan := bc.dataProvider.GetTaxBracketsAsync(year, province.BC)
	bcTax := entity.BcTax{Year: year}
	bpaEWrap := <-bpaEWrapChan
	if bpaEWrap.Error != nil {
		return bcTax, bpaEWrap.Error
	}
	bcTax.BPA = bpaEWrap.Value
	cppBasicEWrap := <-cppBasicEWrapChan
	if cppBasicEWrap.Error != nil {
		return bcTax, cppBasicEWrap.Error
	}
	bcTax.CPPBasic = cppBasicEWrap.Value
	cppFirstEWrap := <-cppFirstEWrapChan
	if cppFirstEWrap.Error != nil {
		return bcTax, cppFirstEWrap.Error
	}
	bcTax.CPPFirst = cppFirstEWrap.Value
	cppSecondEWrap := <-cppSecondEWrapChan
	if cppSecondEWrap.Error != nil {
		return bcTax, cppSecondEWrap.Error
	}
	bcTax.CPPSecond = cppSecondEWrap.Value
	eipEWrap := <-eipEWrapChan
	if eipEWrap.Error != nil {
		return bcTax, eipEWrap.Error
	}
	bcTax.EIP = eipEWrap.Value
	bracketsEWrap := <-bracketsEWrapChan
	if bracketsEWrap.Error != nil {
		return bcTax, bracketsEWrap.Error
	}
	bcTax.Brackets = bracketsEWrap.Value
	return bcTax, nil
}
