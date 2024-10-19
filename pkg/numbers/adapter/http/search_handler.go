package http

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-openapi/swag"
	"github.com/labstack/echo"

	"github.com/jkmolczan/srv-search/pkg/numbers"
	"github.com/jkmolczan/srv-search/pkg/numbers/adapter/http/models"
)

const (
	searchNumbersPathGroup = "search/numbers"
	searchIndexPath        = "/index/:number"

	numberPathParam = "number"

	// for now approximation is hardcoded, in the future it can be a query param
	approximationValue = 0.1
)

var (
	errNumberNotFound         = fmt.Errorf("number not found")
	errInternalServerError    = fmt.Errorf("internal server error")
	errInvalidNumberPathParam = fmt.Errorf("invalid number path param")
)

//go:generate moq -out indexSearcherMock_test.go . indexSearcher
type indexSearcher interface {
	SearchIndex(value int, approximation float64) (numbers.SearchResult, error)
}

func SetSearchNumbersRoutes(e *echo.Echo, h *SearchHandler) {
	g := e.Group(searchNumbersPathGroup)
	g.GET(searchIndexPath, h.SearchNumberIndex)
}

type SearchHandler struct {
	indexSearcher indexSearcher
	logger        echo.Logger
}

func NewSearchHandler(indexSearcher indexSearcher, logger echo.Logger) *SearchHandler {
	return &SearchHandler{
		indexSearcher: indexSearcher,
		logger:        logger,
	}
}

func (h *SearchHandler) SearchNumberIndex(c echo.Context) error {
	inputNumber, err := strconv.Atoi(c.Param(numberPathParam))
	if err != nil {
		h.logger.Debugf("SearchNumberIndex: invalid number path param: %s", c.Param(numberPathParam))
		return errInvalidNumberPathParam
	}

	result, err := h.indexSearcher.SearchIndex(inputNumber, approximationValue)
	if err != nil {
		if errors.Is(err, numbers.ErrIndexNotFound) {
			h.logger.Debugf("SearchNumberIndex: number not found: %d", inputNumber)
			return errNumberNotFound
		}
		h.logger.Errorf("SearchNumberIndex: error searching number index: %+v", err)
		return errInternalServerError
	}

	resp := models.IndexResponse{
		Index:   swag.Int64(int64(result.Index)),
		Message: result.Message,
		Number:  swag.Int64(int64(result.Value)),
	}

	return c.JSON(http.StatusOK, resp)
}
