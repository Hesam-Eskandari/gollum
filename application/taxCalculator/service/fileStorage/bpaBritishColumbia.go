package fileStorage

import (
	"context"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/bpa/entity"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/errorWrap"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/service/fileStorage/internal/models"
)

func (dp *fileStorage) GetBritishColumbiaBPAAsync(ctx context.Context, year int) <-chan errorWrap.ErrorWrap[entity.BritishColumbiaBPA] {
	resChan := make(chan errorWrap.ErrorWrap[entity.BritishColumbiaBPA], 1)
	go func() {
		defer close(resChan)
		res := errorWrap.ErrorWrap[entity.BritishColumbiaBPA]{}
		defer func() { resChan <- res }()
		if dp.data == nil {
			if err := dp.readJson(); err != nil {
				res.Error = err
				return
			}
		}
		var dataForYear *models.BpaBc
		for _, bpa := range dp.data.BpaBc {
			if bpa.Year == year {
				dataForYear = &bpa
				break
			}
		}
		if dataForYear == nil {
			res.Error = ErrBritishColumbiaBPANotFound
			return
		}
		res.Value = dataForYear.ToEntity()
	}()
	return resChan
}
