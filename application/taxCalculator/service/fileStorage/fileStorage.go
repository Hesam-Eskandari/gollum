package fileStorage

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/bpa/entity"
	entity2 "github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/brackets/entity"
	entity3 "github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/cea/entity"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/constants/province"
	entity4 "github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/cpp/entity"
	entity5 "github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/eip/entity"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/errorWrap"
)

const filePath = "application/taxCalculator/service/fileStorage/tax_data.json"

var ErrFederalBPANotFound = errors.New("federal bpa data not found")
var ErrFederalCEANotFound = errors.New("federal cea data not found")
var ErrFederalCPPNotFound = errors.New("federal cpp data not found")
var ErrFederalEIPNotFound = errors.New("federal eip data not found")
var ErrTaxBracketsNotFound = errors.New("federal tax bracket data not found")
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

func (dp *fileStorage) GetFederalBPAAsync(year int) <-chan errorWrap.ErrorWrap[entity.FederalBPA] {
	resChan := make(chan errorWrap.ErrorWrap[entity.FederalBPA], 1)
	go func() {
		defer close(resChan)
		res := errorWrap.ErrorWrap[entity.FederalBPA]{}
		defer func() { resChan <- res }()
		if dp.data == nil {
			if err := dp.readJson(); err != nil {
				res.Error = err
				return
			}
		}
		var dataForYear *bpaFederal
		for _, bpa := range dp.data.BpaFederal {
			if bpa.Year == year {
				dataForYear = &bpa
				break
			}
		}
		if dataForYear == nil {
			res.Error = ErrFederalBPANotFound
			return
		}
		res.Value = dataForYear.toEntity()
	}()
	return resChan
}

func (dp *fileStorage) GetFederalCEAAsync(year int) <-chan errorWrap.ErrorWrap[entity3.FederalCEA] {
	resChan := make(chan errorWrap.ErrorWrap[entity3.FederalCEA], 1)
	go func() {
		defer close(resChan)
		res := errorWrap.ErrorWrap[entity3.FederalCEA]{}
		defer func() { resChan <- res }()
		if dp.data == nil {
			if err := dp.readJson(); err != nil {
				res.Error = err
				return
			}
		}
		var dataForYear *ceaFederal
		for _, cea := range dp.data.CEA {
			if cea.Year == year {
				dataForYear = &cea
				break
			}
		}
		if dataForYear == nil {
			res.Error = ErrFederalCEANotFound
			return
		}
		res.Value = dataForYear.toEntity()
	}()
	return resChan
}

func (dp *fileStorage) GetFederalCPPAsync(year int) <-chan errorWrap.ErrorWrap[entity4.FederalCPP] {
	resChan := make(chan errorWrap.ErrorWrap[entity4.FederalCPP], 1)
	go func() {
		defer close(resChan)
		res := errorWrap.ErrorWrap[entity4.FederalCPP]{}
		defer func() { resChan <- res }()
		if dp.data == nil {
			if err := dp.readJson(); err != nil {
				res.Error = err
				return
			}
		}
		var dataForYear *cppFederal
		for _, cpp := range dp.data.CPP {
			if cpp.Year == year {
				dataForYear = &cpp
				break
			}
		}
		if dataForYear == nil {
			res.Error = ErrFederalCPPNotFound
			return
		}
		res.Value = dataForYear.toEntity()
	}()
	return resChan
}

func (dp *fileStorage) GetFederalEIPAsync(year int) <-chan errorWrap.ErrorWrap[entity5.FederalEIP] {
	resChan := make(chan errorWrap.ErrorWrap[entity5.FederalEIP], 1)
	go func() {
		defer close(resChan)
		res := errorWrap.ErrorWrap[entity5.FederalEIP]{}
		defer func() { resChan <- res }()
		if dp.data == nil {
			if err := dp.readJson(); err != nil {
				res.Error = err
				return
			}
		}
		var dataForYear *eipFederal
		for _, eip := range dp.data.EIP {
			if eip.Year == year {
				dataForYear = &eip
				break
			}
		}
		if dataForYear == nil {
			res.Error = ErrFederalEIPNotFound
			return
		}
		res.Value = dataForYear.toEntity()
	}()
	return resChan
}

func (dp *fileStorage) GetTaxBracketsAsync(year int, pr province.Province) <-chan errorWrap.ErrorWrap[[]entity2.TaxBracket] {
	resChan := make(chan errorWrap.ErrorWrap[[]entity2.TaxBracket], 1)
	go func() {
		defer close(resChan)
		res := errorWrap.ErrorWrap[[]entity2.TaxBracket]{}
		defer func() { resChan <- res }()
		if dp.data == nil {
			if err := dp.readJson(); err != nil {
				res.Error = err
				return
			}
		}
		dataForYear := make([]bracket, 0, 10)
		brks, err := dp.data.getBracketsForProvince(pr)
		if err != nil {
			res.Error = err
			return
		}
		for _, br := range brks {
			if br.Year == year {
				dataForYear = append(dataForYear, br)
			}
		}
		if len(dataForYear) == 0 {
			res.Error = ErrTaxBracketsNotFound
			return
		}
		res.Value = mapBracketsToEntity(dataForYear)
	}()
	return resChan
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
	EIP        []eipFederal `json:"eip"`
	CPP        []cppFederal `json:"cpp"`
	CEA        []ceaFederal `json:"cea"`
	BpaFederal []bpaFederal `json:"bpaFederal"`
	BpaBc      []struct {
		Year  int     `json:"year"`
		Value float64 `json:"value"`
	} `json:"bpaBc"`
	BracketsFederal []bracket `json:"bracketsFederal"`
	BracketsBc      []bracket `json:"bracketsBc"`
	RRSP            []struct {
		Year                  int     `json:"year"`
		Rate                  float64 `json:"rate"`
		MaxContribution       float64 `json:"maxContribution"`
		OverContributionLimit float64 `json:"overContributionLimit"`
		OverContributionRate  float64 `json:"overContributionRate"`
	} `json:"rrsp"`
}

func (td *taxData) getBracketsForProvince(pr province.Province) ([]bracket, error) {
	switch pr {
	case province.BC:
		return td.BracketsBc, nil
	case province.Federal:
		return td.BracketsFederal, nil
	default:
		return nil, ErrProvinceNotFound

	}
}
