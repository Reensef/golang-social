package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Reensef/golang-social/internal/store"
	"github.com/go-chi/chi/v5"
)

type PostPayload struct {
	Title   *string   `json:"title"`
	Content *string   `json:"content"`
	Tags    *[]string `json:"tags"`
}

type postKey string

const postCtxKey postKey = "post"

func parseCreatedPostPayload(src *store.Post, payload *PostPayload) error {
	var errorsList []string

	// Required
	if payload.Title == nil || *payload.Title == "" {
		errorsList = append(errorsList, "title is required")
	}
	if payload.Title == nil || len(*payload.Title) > 100 {
		errorsList = append(errorsList, "title is too long")
	}
	if payload.Content == nil || *payload.Content == "" {
		errorsList = append(errorsList, "content is required")
	}
	if payload.Content == nil || len(*payload.Content) > 1000 {
		errorsList = append(errorsList, "content is too long")
	}

	if len(errorsList) > 0 {
		return errors.New("Invalid payload: " + strings.Join(errorsList, ", "))
	}

	src.Title = *payload.Title
	src.Content = *payload.Content

	// Not required
	if payload.Tags != nil {
		src.Tags = *payload.Tags
	}

	return nil
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload PostPayload
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	post := store.Post{
		UserID: 1,
	}

	err := parseCreatedPostPayload(&post, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.store.Posts.Create(r.Context(), &post); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if err := jsonDataResponse(w, http.StatusOK, post); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	post.Comments = comments

	if err := jsonDataResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "post_id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	err = app.store.Posts.DeleteByID(r.Context(), id)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrDataNotFound):
			app.resourceNotFoundResponse(w, r, err)
		default:
			app.internalServerErrorResponse(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseUpdatedPostPayload(src *store.Post, payload *PostPayload) error {
	var errorsList []string

	if payload.Title != nil {
		if len(*payload.Title) > 100 {
			errorsList = append(errorsList, "title is too long")
		} else {
			src.Title = *payload.Title
		}
	}
	if payload.Content != nil {
		if len(*payload.Content) > 1000 {
			errorsList = append(errorsList, "content is too long")
		} else {
			src.Content = *payload.Content
		}
	}

	if len(errorsList) > 0 {
		return errors.New("Invalid payload: " + strings.Join(errorsList, ", "))
	}

	if payload.Tags != nil {
		src.Tags = *payload.Tags
	}

	return nil
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	var payload PostPayload
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := parseUpdatedPostPayload(post, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}
	if payload.Content != nil {
		post.Content = *payload.Content
	}

	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		switch {
		case errors.Is(err, store.ErrDataNotFound):
			app.resourceConflictResponse(w, r, err)
		default:
			app.internalServerErrorResponse(w, r, err)
		}
		return
	}

	if err := jsonDataResponse(w, http.StatusOK, post); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}

// Сомнительная история
func (app *application) postContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "post_id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalServerErrorResponse(w, r, err)
			return
		}

		post, err := app.store.Posts.GetByID(r.Context(), id)

		if errors.Is(err, store.ErrDataNotFound) {
			app.resourceNotFoundResponse(w, r, err)
			return
		} else if err != nil {
			app.internalServerErrorResponse(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), postCtxKey, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postCtxKey).(*store.Post)
	return post
}
