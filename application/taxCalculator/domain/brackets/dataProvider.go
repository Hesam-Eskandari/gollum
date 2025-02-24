package brackets

import (
	"context"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/brackets/entity"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/constants/province"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/errorWrap"
)

type TaxBracketDataProvider interface {
	GetTaxBracketsAsync(ctx context.Context, year int, pr province.Province) <-chan errorWrap.ErrorWrap[[]entity.TaxBracket]
}
