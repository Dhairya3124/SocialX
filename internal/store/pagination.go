package store

import (
	"net/http"
	"strconv"
)

type PaginatedFeedQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=20"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
	Title  string  `json:"title"`
	Content string `json:"content"`
	Tags  []string `json:"tags"`

}

func (fq PaginatedFeedQuery) Pagination(r *http.Request) (PaginatedFeedQuery, error) {
	limit := r.URL.Query().Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, err
		}
		fq.Limit = l
	}
	offset := r.URL.Query().Get("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return fq, err
		}
		fq.Offset = o
	}
	sort := r.URL.Query().Get("sort")
	if sort != "" {
		fq.Sort = sort
	}

	return fq, nil

}
