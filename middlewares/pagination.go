package middlewares

import (
	"context"
	"net/http"
	"strconv"
)

type contextKey string

type Pagination struct {
	Page  int
	Limit int
}

const (
	minLimit                 = 3
	maxLimit                 = 30
	defaultLimit             = 9
	paginationKey contextKey = "pagination"
)

func PaginationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("page")
		limitStr := r.URL.Query().Get("limit")
		
		page := pageHandler(pageStr)
		limit := limitHandler(limitStr)
		
		ctx := context.WithValue(r.Context(), paginationKey, &Pagination{Page: page, Limit: limit})
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func limitHandler(limitAsString string) int {
	limit, err := strconv.Atoi(limitAsString)
	if err != nil {
		return defaultLimit
	}
	if limit < minLimit {
		return minLimit
	}
	if limit > maxLimit {
		return maxLimit
	}

	return limit
}

func pageHandler(pageAsString string) int {
	page, err := strconv.Atoi(pageAsString)
	
	if err != nil || page < 1 {
		return 1
	}
	
	return page
}

func GetPagination(r *http.Request) Pagination {
	pagination, ok := r.Context().Value(paginationKey).(*Pagination)
	
	if !ok {
		return Pagination{
			Page: 1,
			Limit: defaultLimit,
		}
	}

	return *pagination
}

func CalculateOffset(pagination Pagination) int {
	offset := (pagination.Limit * pagination.Page) - pagination.Limit
	return offset
}