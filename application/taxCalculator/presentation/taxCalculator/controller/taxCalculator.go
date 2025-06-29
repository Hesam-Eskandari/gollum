package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/constants/province"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/tax/interactor"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/presentation/internal/routes"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/service/fileStorage"
	"github.com/Hesam-Eskandari/gollum/library/httpServer"
	"github.com/Hesam-Eskandari/gollum/library/httpServer/middleware"
)

func NewTaxCalculatorController() httpServer.Controller {
	return &taxCalculatorImpl{
		methodValidator: middleware.NewHttpMethodValidator(),
	}
}

type taxCalculatorImpl struct {
	methodValidator middleware.HttpMethodValidator
}

func (tc *taxCalculatorImpl) GetOrderedMiddlewares() []httpServer.Middleware {
	tc.methodValidator.SetAllowedMethods([]string{http.MethodGet})
	return []httpServer.Middleware{tc.methodValidator}
}

func (tc *taxCalculatorImpl) GetUrl() string {
	return routes.TaxCalculatorUrl
}

func (tc *taxCalculatorImpl) Handle(writer http.ResponseWriter, req *http.Request) {
	qParams := req.URL.Query()
	if !qParams.Has("income") {
		http.Error(writer, "income is not provided", http.StatusBadRequest)
		return
	}
	income, err := strconv.ParseFloat(strings.TrimSpace(qParams.Get("income")), 64)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	rrsp := 0.0
	if qParams.Has("rrsp") {
		rrsp, err = strconv.ParseFloat(strings.TrimSpace(qParams.Get("rrsp")), 64)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
	}

	fhsa := 0.0
	if qParams.Has("fhsa") {
		fhsa, err = strconv.ParseFloat(strings.TrimSpace(qParams.Get("fhsa")), 64)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if !qParams.Has("year") {
		http.Error(writer, "year is not provided", http.StatusBadRequest)
		return
	}
	year, err := strconv.Atoi(strings.TrimSpace(qParams.Get("year")))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if !qParams.Has("province") {
		http.Error(writer, "province is not provided", http.StatusBadRequest)
		return
	}
	prStr := strings.TrimSpace(qParams.Get("province"))
	pr, err := province.FromString(prStr)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if pr != province.BC {
		http.Error(writer, "province is not supported", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(req.Context(), time.Second*3)
	defer cancel()
	dataProvider := fileStorage.GetFileStorageInstance()
	federalTaxCalculator := interactor.NewFederalTax(dataProvider)
	fedTax, err := federalTaxCalculator.Calculate(ctx, year, income-fhsa-rrsp)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	fedTaxRaw, err := federalTaxCalculator.Calculate(ctx, year, income)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	bcTaxInteractor := interactor.NewBcTax(dataProvider)
	bcTax, err := bcTaxInteractor.Calculate(ctx, year, income-fhsa-rrsp)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	bcTaxRaw, err := bcTaxInteractor.Calculate(ctx, year, income)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	resp := responseModel{
		Year:                    year,
		Province:                prStr,
		FederalBpa:              fedTax.BPA,
		BcBpa:                   bcTax.BPA,
		Cea:                     fedTax.CEA,
		Eip:                     fedTax.EIP,
		CppBasic:                fedTax.CPPBasic,
		CppFirstAdditional:      fedTax.CPPFirst,
		CppSecondAdditional:     fedTax.CPPSecond,
		PayableFederalTax:       fedTax.PayableTax,
		TotalFederalTax:         fedTax.TotalTax,
		PayableBcTax:            bcTax.PayableTax,
		TotalBcTax:              bcTax.TotalTax,
		TotalFederalCredits:     fedTax.TotalCredits,
		TotalBcCredits:          bcTax.TotalCredits,
		TotalFederalDeductions:  float64(int((fedTax.TotalDeductions+rrsp+fhsa)*1000)) / 1000,
		TotalBcDeductions:       float64(int((bcTax.TotalDeductions+rrsp+fhsa)*1000)) / 1000,
		FederalCreditsReduction: fedTax.CreditsTaxReduction,
		BcCreditsReduction:      bcTax.CreditsTaxReduction,
		AfterTaxIncome:          float64(int((income-fedTax.PayableTax-bcTax.PayableTax-fedTax.EIP-fedTax.CPPBasic-fedTax.CPPFirst-fedTax.CPPSecond)*1000)) / 1000,
		TaxReturn:               float64(int((fedTaxRaw.PayableTax-fedTax.PayableTax+bcTaxRaw.PayableTax-bcTax.PayableTax)*1000)) / 1000,
		RRSPContribution:        rrsp,
		FHSAContribution:        fhsa,
	}
	body, err := json.Marshal(&resp)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err = writer.Write(body); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
}

type responseModel struct {
	Year                    int     `json:"year"`
	Province                string  `json:"province"`
	FederalBpa              float64 `json:"federalBpa"`
	BcBpa                   float64 `json:"bcBpa"`
	Cea                     float64 `json:"cea"`
	Eip                     float64 `json:"eip"`
	CppBasic                float64 `json:"cppBasic"`
	CppFirstAdditional      float64 `json:"cppFirstAdditional"`
	CppSecondAdditional     float64 `json:"cppSecondAdditional"`
	PayableFederalTax       float64 `json:"payableFederalTax"`
	TotalFederalTax         float64 `json:"totalFederalTax"`
	PayableBcTax            float64 `json:"payableBcTax"`
	TotalBcTax              float64 `json:"totalBcTax"`
	TotalFederalCredits     float64 `json:"totalFederalCredits"`
	TotalBcCredits          float64 `json:"totalBcCredits"`
	TotalFederalDeductions  float64 `json:"totalFederalDeductions"`
	TotalBcDeductions       float64 `json:"totalBcDeductions"`
	FederalCreditsReduction float64 `json:"federalCreditsReduction"`
	BcCreditsReduction      float64 `json:"bcCreditsReduction"`
	AfterTaxIncome          float64 `json:"afterTaxIncome"`
	TaxReturn               float64 `json:"taxReturn"`
	RRSPContribution        float64 `json:"rrspContribution"`
	FHSAContribution        float64 `json:"fhsaContribution"`
}
