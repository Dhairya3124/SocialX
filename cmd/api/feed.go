package main

import (
	"net/http"

	"github.com/Dhairya3124/SocialX/internal/store"
)

func (app *application) userFeedHandler(w http.ResponseWriter, r *http.Request) {
	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}
	fq, err := fq.Pagination(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	ctx := r.Context()
	feed, err := app.store.Posts.GetUserFeed(ctx, int64(8), fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}

}
