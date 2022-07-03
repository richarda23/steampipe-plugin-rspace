package rspace

import (
	"fmt"

	"github.com/richarda23/rspace-client-go/rspace"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func listingHasNextPage(links []rspace.Link) bool {
	for _, v := range links {

		if v.Rel == "next" {
			return true
		}
	}
	return false
}

// calculatePageSizes calculates a list of pagesizes to use given a SQL limit
func calculatePageSizes(lim int, HARD_LIMIT int, maxPageSize int) ([]int, error) {

	if lim < 1 || HARD_LIMIT < 1 || maxPageSize < 1 {
		return nil, fmt.Errorf("invalid arguments: must be positive integers > 1 - %d, %d ,%d", lim, HARD_LIMIT, maxPageSize)
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

// getLimit extracts the  limit specifed in the database query, or returns HARD_LIMIT
func getLimit(d *plugin.QueryData) int {
	limit := HARD_LIMIT
	if d != nil && d.QueryContext != nil && d.QueryContext.Limit != nil {
		lim := d.QueryContext.Limit
		if *lim > 0 && *lim < HARD_LIMIT {
			limit = int(*lim)
		}
	}
	return limit
}
