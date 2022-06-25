package rspace

import (
	"testing"
)

var hardLimit = 1000
var maxPageSize = 20

func TestDBLimitLessThanPageSize(t *testing.T) {
	dblimit := 11
	expected := 11
	res, _ := calculatePageSizes(int64(dblimit), hardLimit, maxPageSize)
	if len(res) != 1 && res[0] != expected {
		t.Errorf("expected %d but was %d", expected, res[0])
	}
}

func TestDBLimitGreaterThanHardLimit(t *testing.T) {
	dblimit := hardLimit + 1
	maxPageSize = 200
	expected := []int{maxPageSize, maxPageSize, maxPageSize, maxPageSize, maxPageSize}
	res, _ := calculatePageSizes(int64(dblimit), hardLimit, maxPageSize)
	if len(res) != len(expected) {
		t.Errorf("expected %d but was %d", len(res), len(expected))
	}
	for _, v := range res {
		if v != maxPageSize {
			t.Errorf("expected %d but was %d", maxPageSize, v)
		}
	}
}
