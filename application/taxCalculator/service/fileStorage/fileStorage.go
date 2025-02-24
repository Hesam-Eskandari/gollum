package fileStorage

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/constants/province"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/service/fileStorage/internal/models"
)

const filePath = "application/taxCalculator/service/fileStorage/data/tax_data.json"

var ErrFederalBPANotFound = errors.New("federal bpa data not found")
var ErrBritishColumbiaBPANotFound = errors.New("british columbia bpa data not found")
var ErrFederalCEANotFound = errors.New("federal cea data not found")
var ErrFederalCPPNotFound = errors.New("federal cpp data not found")
var ErrFederalEIPNotFound = errors.New("federal eip data not found")
var ErrTaxBracketsNotFound = errors.New("tax Bracket data not found")
var ErrProvinceNotFound = errors.New("province not found")

var instance *fileStorage

func GetFileStorageInstance() domain.DataProvider {
	if instance != nil {
		return instance
	}
	instance = &fileStorage{}
	go instance.readJson()
	return instance
}

type fileStorage struct {
	data *taxData
}

func (dp *fileStorage) readJson() error {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 400)
	if err != nil {
		return err
	}
	defer file.Close()
	var data taxData
	if err = json.NewDecoder(file).Decode(&data); err != nil {
		return err
	}
	dp.data = &data
	return nil
}

type taxData struct {
	EIP             []models.EipFederal `json:"eip"`
	CPP             []models.CppFederal `json:"cpp"`
	CEA             []models.CeaFederal `json:"cea"`
	BpaFederal      []models.BpaFederal `json:"bpaFederal"`
	BpaBc           []models.BpaBc      `json:"bpaBc"`
	BracketsFederal []models.Bracket    `json:"bracketsFederal"`
	BracketsBc      []models.Bracket    `json:"bracketsBc"`
	RRSP            []struct {
		Year                  int     `json:"year"`
		Rate                  float64 `json:"rate"`
		MaxContribution       float64 `json:"maxContribution"`
		OverContributionLimit float64 `json:"overContributionLimit"`
		OverContributionRate  float64 `json:"overContributionRate"`
	} `json:"rrsp"`
}

func (td *taxData) getBracketsForProvince(pr province.Province) ([]models.Bracket, error) {
	switch pr {
	case province.BC:
		return td.BracketsBc, nil
	case province.Federal:
		return td.BracketsFederal, nil
	default:
		return nil, ErrProvinceNotFound

	}
}
