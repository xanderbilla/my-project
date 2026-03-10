package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/xanderbilla/my-project/internal/store"
)

type CreatePostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload

	if err := DecodeJSON(w, r, &payload); err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		// TODO: change after auth
		UserID: 1,
	}
	ctx := r.Context()

	if err := app.storage.Posts.Create(ctx, post); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := Success(w, http.StatusCreated, "Post created successfully", post); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postId")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()

	post, err := app.storage.Posts.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			Error(w, http.StatusNotFound, err.Error())
		default:
			Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if err := Success(w, http.StatusOK, "Post fetched successfully", post); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
}
