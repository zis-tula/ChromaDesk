package handler

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaginationParams struct {
	Page   int
	Size   int
	Sort   string
	Order  string
	Search string
}

func ParsePagination(c *gin.Context) PaginationParams {
	p := PaginationParams{
		Page:   1,
		Size:   20,
		Sort:   "created_at",
		Order:  "desc",
		Search: c.Query("search"),
	}

	if page, err := strconv.Atoi(c.Query("page")); err == nil && page > 0 {
		p.Page = page
	}
	if size, err := strconv.Atoi(c.Query("size")); err == nil && size > 0 {
		if size > 100 {
			size = 100
		}
		p.Size = size
	}
	if sort := c.Query("sort"); sort != "" {
		p.Sort = sort
	}
	if order := c.Query("order"); order == "asc" || order == "desc" {
		p.Order = order
	}

	return p
}

func (p PaginationParams) Offset() int {
	return (p.Page - 1) * p.Size
}

func (p PaginationParams) OrderClause() string {
	return p.Sort + " " + p.Order
}

type PaginatedResponse struct {
	Items      interface{} `json:"items"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Size       int         `json:"size"`
	TotalPages int         `json:"totalPages"`
}

func NewPaginatedResponse(items interface{}, total int64, p PaginationParams) PaginatedResponse {
	totalPages := int(math.Ceil(float64(total) / float64(p.Size)))
	return PaginatedResponse{
		Items:      items,
		Total:      total,
		Page:       p.Page,
		Size:       p.Size,
		TotalPages: totalPages,
	}
}
