package models

type PaginationParams struct {
	Page     int `query:"page"`
	PageSize int `query:"page_size"`
}