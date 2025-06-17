package response

// CalculatePagination calculates pagination metadata
func CalculatePagination(params PaginationParams, total int) *Meta {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PerPage < 1 {
		params.PerPage = 10 // Default items per page
	}

	totalPages := (total + params.PerPage - 1) / params.PerPage

	return &Meta{
		Page:      params.Page,
		PerPage:   params.PerPage,
		Total:     total,
		TotalPage: totalPages,
	}
}

// ParsePaginationParams extracts and validates pagination parameters
func ParsePaginationParams(page, perPage int) PaginationParams {
	params := PaginationParams{
		Page:    page,
		PerPage: perPage,
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.PerPage < 1 || params.PerPage > 100 {
		params.PerPage = 10 // Default to 10 items per page with a max of 100
	}

	return params
}

// GetOffset calculates the database query offset based on pagination params
func (p PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.PerPage
}

// GetLimit returns the number of items to fetch
func (p PaginationParams) GetLimit() int {
	return p.PerPage
}
