package numbers

import (
	"fmt"
	"math"
)

const (
	msgValueFound            = "Value found."
	msgApproximateValueFound = "Approximate value found."
)

var (
	ErrIndexNotFound           = fmt.Errorf("index not found")
	ErrApproximationOutOfRange = fmt.Errorf("approximation out of range")

	approximationMin = 0.0
	approximationMax = 0.1
)

type SearchResult struct {
	Index   int
	Value   int
	Message string
}

//go:generate moq -out numbersStorageMock_test.go . numbersStorage
type numbersStorage interface {
	GetSortedCollection() ([]int, error)
}

type SearchService struct {
	numbersStorage numbersStorage
}

func NewSearchService(numbersStorage numbersStorage) *SearchService {
	return &SearchService{numbersStorage: numbersStorage}
}

// SearchIndex searches index of the provided value in the sorted collection with the given approximation using binary search.
func (ss *SearchService) SearchIndex(value int, approximation float64) (SearchResult, error) {
	if approximation < approximationMin || approximation > approximationMax {
		return SearchResult{}, ErrApproximationOutOfRange
	}
	data, err := ss.numbersStorage.GetSortedCollection()
	if err != nil {
		return SearchResult{}, fmt.Errorf("failed to get sorted collection: %w", err)
	}

	left, right := 0, len(data)-1
	for left <= right {
		mid := left + (right-left)/2
		if data[mid] == value {
			return SearchResult{
				Index:   mid,
				Value:   data[mid],
				Message: msgValueFound,
			}, nil
		}
		if data[mid] < value {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	// approximate match within provided approximation
	var candidates []int
	if left < len(data) {
		candidates = append(candidates, left)
	}
	if right >= 0 {
		candidates = append(candidates, right)
	}

	for _, idx := range candidates {
		if ss.withinApproximation(data[idx], value, approximation) {
			return SearchResult{
				Index:   idx,
				Value:   data[idx],
				Message: msgApproximateValueFound,
			}, nil
		}
	}

	return SearchResult{}, ErrIndexNotFound
}

func (ss *SearchService) withinApproximation(a, b int, approximation float64) bool {
	if b == 0 {
		return a == 0
	}
	diff := math.Abs(float64(a - b))
	return (diff / float64(b)) <= approximation
}
