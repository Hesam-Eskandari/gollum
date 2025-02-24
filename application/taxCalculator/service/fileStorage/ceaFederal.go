package fileStorage

import (
	"context"
	entity3 "github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/cea/entity"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/errorWrap"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/service/fileStorage/internal/models"
)

func (dp *fileStorage) GetFederalCEAAsync(ctx context.Context, year int) <-chan errorWrap.ErrorWrap[entity3.FederalCEA] {
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
		var dataForYear *models.CeaFederal
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
		res.Value = dataForYear.ToEntity()
	}()
	return resChan
}
