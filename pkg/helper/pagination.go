package helper

import (
	"context"
)

type Pagination struct {
	Limit     int `json:"limit_per_page"`
	Page      int `json:"current_page"`
	TotalPage int `json:"total_page"`
	TotalRows int `json:"total_rows"`
}

func CalculatePagination(ctx context.Context, limit, totalRows int) (*Pagination, error) {

	// Calculate total pages
	totalPages := totalRows / limit
	if totalRows%limit != 0 {
		totalPages++
	}

	return &Pagination{
		Limit:     limit,
		TotalRows: totalRows,
		TotalPage: totalPages,
	}, nil
}
