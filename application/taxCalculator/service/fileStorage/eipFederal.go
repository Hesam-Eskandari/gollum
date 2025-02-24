package fileStorage

import (
	entity5 "github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/eip/entity"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/errorWrap"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/service/fileStorage/internal/models"
)

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
		var dataForYear *models.EipFederal
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
		res.Value = dataForYear.ToEntity()
	}()
	return resChan
}
