package pagination

import (
	"math"
)

type Pagination struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	Total       int `json:"total"`
	LastPage    int `json:"last_page"`
}

func NewPagination(perPage, currentPage, total int) Pagination {
	lastPage := int(math.Ceil(float64(total) / float64(perPage)))

	if currentPage == -1 {
		lastPage = 1
	}

	return Pagination{
		CurrentPage: currentPage,
		PerPage:     perPage,
		Total:       total,
		LastPage:    lastPage,
	}
}

func (p *Pagination) GetOffset() int {
	return (p.GetCurrentPage() - 1) * p.GetPerPage()
}
func (p *Pagination) GetPerPage() int {
	if p.PerPage == 0 {
		p.PerPage = 7
	}
	return p.PerPage
}
func (p *Pagination) GetCurrentPage() int {
	if p.CurrentPage == 0 {
		p.CurrentPage = 1
	}
	return p.CurrentPage
}
