package entity

type QueryRequest struct {
	Page      int    `query:"page"`
	Limit     int    `query:"limit"`
	SortBy    string `query:"sort_by"`
	Search    string `query:"search"`
	SortOrder string `query:"sort_order"`
}
