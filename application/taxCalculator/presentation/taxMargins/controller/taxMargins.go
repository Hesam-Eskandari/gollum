package controller

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/brackets/interactor"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/domain/constants/province"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/presentation/internal/routes"
	"github.com/Hesam-Eskandari/gollum/application/taxCalculator/service/fileStorage"
	"github.com/Hesam-Eskandari/gollum/library/httpServer"
	"github.com/Hesam-Eskandari/gollum/library/httpServer/middleware"
)

func NewTaxMarginsController() httpServer.Controller {
	return &taxMarginsImpl{
		methodValidator: middleware.NewHttpMethodValidator(),
	}
}

type taxMarginsImpl struct {
	methodValidator middleware.HttpMethodValidator
}

func (tm *taxMarginsImpl) GetOrderedMiddlewares() []httpServer.Middleware {
	tm.methodValidator.SetAllowedMethods([]string{http.MethodGet})
	return []httpServer.Middleware{tm.methodValidator}
}

func (tm *taxMarginsImpl) GetUrl() string {
	return routes.TaxMarginsUrl
}

func (tm *taxMarginsImpl) Handle(writer http.ResponseWriter, req *http.Request) {
	qParams := req.URL.Query()
	if !qParams.Has("province") {
		http.Error(writer, "province parameter not provided", http.StatusBadRequest)
		return
	}
	prStr := strings.TrimSpace(qParams.Get("province"))
	pr, err := province.FromString(prStr)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	year, err := strconv.Atoi(qParams.Get("year"))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if pr != province.BC {
		http.Error(writer, "province is not supported", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(req.Context(), time.Second*3)
	defer cancel()
	dataProvider := fileStorage.GetFileStorageInstance()
	calculator := interactor.NewTaxBracketCalculator(dataProvider)
	taxMargins, err := calculator.GetAllMarginalBrackets(ctx, year, pr)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	respBrackets := make([]RespBracket, 0, 15)
	for _, bracket := range taxMargins {
		respBrackets = append(respBrackets, RespBracket{
			Low:  Float64Ext(bracket.Low),
			High: Float64Ext(bracket.High),
			Rate: Float64Ext(bracket.Rate),
		})
	}
	resp := responseModel{
		Brackets: respBrackets,
	}
	body, err := json.Marshal(&resp)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err = writer.Write(body); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
}

type responseModel struct {
	Brackets []RespBracket `json:"brackets"`
}

type RespBracket struct {
	Low  Float64Ext `json:"low"`
	High Float64Ext `json:"high"`
	Rate Float64Ext `json:"rate"`
}

type Float64Ext float64

func (f Float64Ext) MarshalJSON() ([]byte, error) {
	switch {
	case math.IsInf(float64(f), 1):
		return []byte(`"Infinity"`), nil
	case math.IsInf(float64(f), -1):
		return []byte(`"-Infinity"`), nil
	case math.IsNaN(float64(f)):
		return []byte(`"NaN"`), nil
	default:
		return json.Marshal(float64(f))
	}
}
