package brackets

import (
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/brackets/entity"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/constants/province"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/errorWrap"
)

type TaxBracketDataProvider interface {
	GetTaxBracketsAsync(year int, pr province.Province) <-chan errorWrap.ErrorWrap[[]entity.TaxBracket]
}
