package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"

	"github.com/jkmolczan/srv-search/pkg/numbers/adapter/http/models"
)

func ErrorHandler(err error, c echo.Context) {
	var he *echo.HTTPError
	ok := errors.As(err, &he)
	if ok {
		_ = c.JSON(he.Code, he.Message)
		return
	}

	if errors.Is(err, errNumberNotFound) {
		_ = c.JSON(http.StatusNotFound, &models.Error{Message: err.Error()})
		return
	}

	if errors.Is(err, errInvalidNumberPathParam) {
		_ = c.JSON(http.StatusBadRequest, &models.Error{Message: err.Error()})
		return
	}

	if errors.Is(err, errInternalServerError) {
		_ = c.JSON(http.StatusInternalServerError, &models.Error{Message: err.Error()})
		return
	}

}
