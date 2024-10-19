package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"

	"github.com/jkmolczan/srv-search/pkg/numbers"
	"github.com/jkmolczan/srv-search/pkg/numbers/adapter/http/models"
)

func TestSearchHandler_SearchNumberIndex(t *testing.T) {
	e := echo.New()
	e.HTTPErrorHandler = ErrorHandler
	e.Logger.SetLevel(log.OFF)

	t.Run("should return http 400 error on invalid number path param", func(t *testing.T) {
		h := NewSearchHandler(nil, e.Logger)
		SetSearchNumbersRoutes(e, h)

		request := httptest.NewRequest(http.MethodGet, "/search/numbers/index/123aa", nil)

		requestRecorder := httptest.NewRecorder()
		e.ServeHTTP(requestRecorder, request)

		assert.Equal(t, http.StatusBadRequest, requestRecorder.Code)

		var respBody models.Error
		if err := respBody.UnmarshalBinary(requestRecorder.Body.Bytes()); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "invalid number path param", respBody.Message)
	})

	t.Run("should return http 404 error on number not found", func(t *testing.T) {
		s := &indexSearcherMock{
			SearchIndexFunc: func(value int, approximation float64) (numbers.SearchResult, error) {
				return numbers.SearchResult{}, numbers.ErrIndexNotFound
			},
		}

		h := NewSearchHandler(s, e.Logger)
		SetSearchNumbersRoutes(e, h)

		request := httptest.NewRequest(http.MethodGet, "/search/numbers/index/123", nil)

		requestRecorder := httptest.NewRecorder()
		e.ServeHTTP(requestRecorder, request)

		assert.Equal(t, http.StatusNotFound, requestRecorder.Code)

		var respBody models.Error
		if err := respBody.UnmarshalBinary(requestRecorder.Body.Bytes()); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "number not found", respBody.Message)
	})

	t.Run("should return http 500 error on internal server error", func(t *testing.T) {
		s := &indexSearcherMock{
			SearchIndexFunc: func(value int, approximation float64) (numbers.SearchResult, error) {
				return numbers.SearchResult{}, assert.AnError
			},
		}

		h := NewSearchHandler(s, e.Logger)
		SetSearchNumbersRoutes(e, h)

		request := httptest.NewRequest(http.MethodGet, "/search/numbers/index/123", nil)

		requestRecorder := httptest.NewRecorder()
		e.ServeHTTP(requestRecorder, request)

		assert.Equal(t, http.StatusInternalServerError, requestRecorder.Code)

		var respBody models.Error
		if err := respBody.UnmarshalBinary(requestRecorder.Body.Bytes()); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "internal server error", respBody.Message)
	})

	t.Run("should return http 200 on successful search", func(t *testing.T) {
		s := &indexSearcherMock{
			SearchIndexFunc: func(value int, approximation float64) (numbers.SearchResult, error) {
				assert.Equal(t, 123, value)
				assert.Equal(t, approximationValue, approximation)
				return numbers.SearchResult{
					Index:   32,
					Value:   value,
					Message: "Found index: 32 of the provided number: 123",
				}, nil
			},
		}

		h := NewSearchHandler(s, e.Logger)
		SetSearchNumbersRoutes(e, h)

		request := httptest.NewRequest(http.MethodGet, "/search/numbers/index/123", nil)

		requestRecorder := httptest.NewRecorder()
		e.ServeHTTP(requestRecorder, request)

		assert.Equal(t, http.StatusOK, requestRecorder.Code)

		var respBody models.IndexResponse
		if err := respBody.UnmarshalBinary(requestRecorder.Body.Bytes()); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, int64(32), *respBody.Index)
		assert.Equal(t, "Found index: 32 of the provided number: 123", respBody.Message)
		assert.Equal(t, int64(123), *respBody.Number)
	})
}
