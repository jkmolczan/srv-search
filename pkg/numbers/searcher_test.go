package numbers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearcher_SearchIndex(t *testing.T) {
	tests := map[string]struct {
		numbersCollection []int
		numberToFind      int
		approximation     float64
		expectedIndex     int
		expectedValue     int
		expectedMessage   string
		expectedError     error
	}{
		"should return exact value": {
			numbersCollection: []int{1, 4, 6, 8, 9, 10, 12, 13, 15, 17, 19},
			numberToFind:      13,
			approximation:     0.1,
			expectedIndex:     7,
			expectedValue:     13,
			expectedMessage:   msgValueFound,
			expectedError:     nil,
		},

		"should return exact value for one element collection": {
			numbersCollection: []int{13},
			numberToFind:      13,
			approximation:     0.1,
			expectedIndex:     0,
			expectedValue:     13,
			expectedMessage:   msgValueFound,
			expectedError:     nil,
		},
		"should return exact value (first element)": {
			numbersCollection: []int{1, 4, 6, 8, 9, 10, 12, 13, 15, 17, 19},
			numberToFind:      1,
			approximation:     0.1,
			expectedIndex:     0,
			expectedValue:     1,
			expectedMessage:   msgValueFound,
			expectedError:     nil,
		},
		"should return exact value (last element)": {
			numbersCollection: []int{1, 4, 6, 8, 9, 10, 12, 13, 15, 17, 19},
			numberToFind:      19,
			approximation:     0.1,
			expectedIndex:     10,
			expectedValue:     19,
			expectedMessage:   msgValueFound,
			expectedError:     nil,
		},
		"should return approximate value": {
			numbersCollection: []int{1, 4, 6, 8, 9, 10, 12, 13, 15, 17, 19},
			numberToFind:      18,
			approximation:     0.1,
			expectedIndex:     10,
			expectedValue:     19,
			expectedMessage:   msgApproximateValueFound,
			expectedError:     nil,
		},
		"should return approximate value for one element collection": {
			numbersCollection: []int{57},
			numberToFind:      53,
			approximation:     0.1,
			expectedIndex:     0,
			expectedValue:     57,
			expectedMessage:   msgApproximateValueFound,
			expectedError:     nil,
		},

		"should return approximate value (first element)": {
			numbersCollection: []int{10, 17, 19, 20, 22, 24, 25, 29},
			numberToFind:      11,
			approximation:     0.1,
			expectedIndex:     0,
			expectedValue:     10,
			expectedMessage:   msgApproximateValueFound,
			expectedError:     nil,
		},

		"should return approximate value (last element)": {
			numbersCollection: []int{10, 17, 19, 20, 22, 24, 25, 29},
			numberToFind:      30,
			approximation:     0.1,
			expectedIndex:     7,
			expectedValue:     29,
			expectedMessage:   msgApproximateValueFound,
			expectedError:     nil,
		},

		"should return exact value if approximation is set to 0": {
			numbersCollection: []int{1, 4, 6, 8, 9, 10, 12, 13, 15, 17, 19},
			numberToFind:      17,
			approximation:     0.0,
			expectedIndex:     9,
			expectedValue:     17,
			expectedMessage:   msgValueFound,
			expectedError:     nil,
		},

		"should return error when index not found and approximation is set to 0": {
			numbersCollection: []int{1, 4, 6, 8, 9, 10, 12, 13, 15, 17, 19},
			numberToFind:      18,
			approximation:     0.0,
			expectedIndex:     0,
			expectedValue:     0,
			expectedMessage:   "",
			expectedError:     ErrIndexNotFound,
		},

		"should return error when index not found and collection of the numbers is empty": {
			numbersCollection: []int{},
			numberToFind:      13,
			approximation:     0.1,
			expectedIndex:     0,
			expectedValue:     0,
			expectedMessage:   "",
			expectedError:     ErrIndexNotFound,
		},

		"should return error when index not found and collection of the numbers has one element": {
			numbersCollection: []int{57},
			numberToFind:      13,
			approximation:     0.1,
			expectedIndex:     0,
			expectedValue:     0,
			expectedMessage:   "",
			expectedError:     ErrIndexNotFound,
		},

		"should return error when index not found and collection of the numbers is not empty": {
			numbersCollection: []int{1, 4, 6, 8, 9, 10, 12, 13, 15, 17, 19},
			numberToFind:      3,
			approximation:     0.1,
			expectedIndex:     0,
			expectedValue:     0,
			expectedMessage:   "",
			expectedError:     ErrIndexNotFound,
		},

		"should return error when approximation out of range (too high)": {
			numbersCollection: []int{1, 4, 6, 8, 9, 10, 12, 13, 15, 17, 19},
			numberToFind:      13,
			approximation:     0.2,
			expectedIndex:     0,
			expectedValue:     0,
			expectedMessage:   "",
			expectedError:     ErrApproximationOutOfRange,
		},

		"should return error when approximation out of range (too low)": {
			numbersCollection: []int{1, 4, 6, 8, 9, 10, 12, 13, 15, 17, 19},
			numberToFind:      13,
			approximation:     -0.1,
			expectedIndex:     0,
			expectedValue:     0,
			expectedMessage:   "",
			expectedError:     ErrApproximationOutOfRange,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			storage := &numbersStorageMock{
				GetSortedCollectionFunc: func() ([]int, error) {
					return test.numbersCollection, nil
				},
			}

			searcher := NewSearchService(storage)

			result, err := searcher.SearchIndex(test.numberToFind, test.approximation)
			assert.Equal(t, test.expectedIndex, result.Index)
			assert.Equal(t, test.expectedValue, result.Value)
			assert.Equal(t, test.expectedMessage, result.Message)
			assert.Equal(t, test.expectedError, err)
		})
	}

	t.Run("should return error when failed to get sorted collection", func(t *testing.T) {
		storage := &numbersStorageMock{
			GetSortedCollectionFunc: func() ([]int, error) {
				return nil, assert.AnError
			},
		}

		searcher := NewSearchService(storage)

		result, err := searcher.SearchIndex(13, 0.1)

		assert.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, 0, result.Index)
		assert.Equal(t, 0, result.Value)
		assert.Equal(t, "", result.Message)
	})

}
