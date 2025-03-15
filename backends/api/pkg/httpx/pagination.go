package httpx

import (
	"github.com/danielgtaylor/huma/v2"
)

// SortDirection is the direction of the sort
type SortDirection int

const (
	SortDesc SortDirection = iota
	SortAsc
)

// SQL returns the SQL string for the sort direction
func (s SortDirection) SQL() string {
	if s == SortAsc {
		return "ASC"
	}
	return "DESC"
}

// CursorPaginationParams is the parameters for cursor pagination, only one of After or Before should be provided
type CursorPaginationParams struct {
	After     string
	Before    string
	Direction SortDirection
	PageSize  int
}

// CursorPagination is the pagination for cursor pagination
type CursorPagination struct {
	After     string `query:"after" required:"false" doc:"Pagination cursor for fetching next page"`
	Before    string `query:"before" required:"false" doc:"Pagination cursor for fetching previous page"`
	Direction string `query:"direction" enum:"asc,desc" doc:"Sort direction for items" default:"desc"`
	PageSize  int    `query:"pageSize" doc:"Maximum number of items to return" minimum:"10" maximum:"100" default:"10"`

	Params CursorPaginationParams `hidden:"true"`
}

// Resolve resolves the cursor pagination parameters
func (p *CursorPagination) Resolve(ctx huma.Context, prefix *huma.PathBuffer) []error {
	dir := SortDesc
	if p.Direction == "asc" {
		dir = SortAsc
	}

	pageSize := p.PageSize
	if pageSize <= 0 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100
	}

	if p.After != "" && p.Before != "" {
		return []error{ErrInvalidCursor.WithLocation(prefix.With("cursor"))}
	}

	p.Params = CursorPaginationParams{
		After:     p.After,
		Before:    p.Before,
		Direction: dir,
		PageSize:  pageSize,
	}

	return nil
}
