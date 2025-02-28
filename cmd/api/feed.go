package main

import "net/http"

func (app *application) userFeedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	feed, err := app.store.Posts.GetUserFeed(ctx, int64(7))
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}

}
