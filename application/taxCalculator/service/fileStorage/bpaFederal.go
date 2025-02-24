package fileStorage

import (
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/bpa/entity"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/errorWrap"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/service/fileStorage/internal/models"
)

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
		var dataForYear *models.BpaFederal
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
		res.Value = dataForYear.ToEntity()
	}()
	return resChan
}
