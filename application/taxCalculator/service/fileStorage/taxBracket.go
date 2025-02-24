package fileStorage

import (
	"context"
	entity2 "github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/brackets/entity"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/constants/province"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/errorWrap"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/service/fileStorage/internal/models"
)

func (dp *fileStorage) GetTaxBracketsAsync(ctx context.Context, year int, pr province.Province) <-chan errorWrap.ErrorWrap[[]entity2.TaxBracket] {
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
		dataForYear := make([]models.Bracket, 0, 10)
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
		res.Value = models.MapBracketsToEntity(dataForYear)
	}()
	return resChan
}
