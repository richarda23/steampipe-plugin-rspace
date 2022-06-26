package rspace

import (
	"errors"
	"fmt"

	"github.com/richarda23/rspace-client-go/rspace"
)

func listingHasNextPage(links []rspace.Link) bool {
	for _, v := range links {

		if v.Rel == "next" {
			return true
		}
	}
	return false
}

// calculatePageSizes calculates a list of pagesizes to use given a SQL limit and
func calculatePageSizes(lim int, HARD_LIMIT int, maxPageSize int) ([]int, error) {

	if lim < 1 || HARD_LIMIT < 1 || maxPageSize < 1 {
		return nil, errors.New(fmt.Sprintf("Invalid arguments: must be positive integers > 1 - %d, %d ,%d", lim, HARD_LIMIT, maxPageSize))
	}
	if lim < maxPageSize {
		return []int{int(lim)}, nil
	} else if lim > HARD_LIMIT {
		len := HARD_LIMIT / maxPageSize
		pageSizes := make([]int, len)
		for i := 0; i < len; i += 1 {
			pageSizes[i] = maxPageSize
		}
		return pageSizes, nil
	} else {
		add_remainder_page := 0
		limInt := int(lim)
		remainder := limInt % maxPageSize
		if remainder > 0 {
			add_remainder_page = 1
		}
		num_items := (limInt / maxPageSize) + add_remainder_page
		final_page_size := remainder
		pageSizes := make([]int, num_items)
		for i := 0; i < num_items-add_remainder_page; i += 1 {
			pageSizes[i] = maxPageSize
		}
		if remainder > 0 {
			pageSizes[num_items-1] = final_page_size

		}
		return pageSizes, nil
	}

}
