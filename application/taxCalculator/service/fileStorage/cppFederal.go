package fileStorage

import (
	entity4 "github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/cpp/entity"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/errorWrap"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/service/fileStorage/internal/models"
)

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
		var dataForYear *models.CppFederal
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
		res.Value = dataForYear.ToEntity()
	}()
	return resChan
}
