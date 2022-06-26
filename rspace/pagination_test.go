package rspace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var hardLimit = 1000
var maxPageSize = 20

func TestDBLimitLessThanPageSize(t *testing.T) {
	dblimit := 11
	expected := 11
	res, _ := calculatePageSizes(dblimit, hardLimit, maxPageSize)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, expected, res[0])
}

func TestDBLimitGreaterThanHardLimit(t *testing.T) {
	dblimit := hardLimit + 1
	maxPageSize = 200
	expected := []int{maxPageSize, maxPageSize, maxPageSize, maxPageSize, maxPageSize}
	res, _ := calculatePageSizes(dblimit, hardLimit, maxPageSize)
	assert.Equal(t, len(expected), len(res))
	assert.ElementsMatch(t, expected, res)
}

type paginationInput struct {
	dbLimit, maxPageSize int
	expected             []int
}

func TestDBLimitInBetween(t *testing.T) {
	inputs := []paginationInput{
		{201, 200, []int{200, 1}},
		{999, 200, []int{200, 200, 200, 200, 199}},
		{3, 2, []int{2, 1}},
		{3, 1, []int{1, 1, 1}},
		{400, 200, []int{200, 200}},
	}

	for i, input := range inputs {
		res, _ := calculatePageSizes(input.dbLimit, hardLimit, input.maxPageSize)
		assert.Equal(t, len(input.expected), len(res), "input %d failed", i)
		assert.ElementsMatch(t, input.expected, res, "input %d failed", i)
	}

}

func TestDBLimitErrorHandling(t *testing.T) {
	input := paginationInput{0, 200, nil}

	res, _ := calculatePageSizes(input.dbLimit, hardLimit, input.maxPageSize)
	assert.Nil(t, res)

}
