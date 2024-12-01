package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/Reensef/golang-social/internal/store"
	"github.com/go-chi/chi/v5"
)

type CommentPayload struct {
	UserID  *int64  `json:"user_id"`
	Content *string `json:"content"`
}

func (app *application) createPostCommentHandler(w http.ResponseWriter, r *http.Request) {
	postIDParam := chi.URLParam(r, "post_id")
	postID, err := strconv.ParseInt(postIDParam, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, ErrInvalidPostID)
		return
	}
	var payload CommentPayload
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	var comment store.Comment
	// TODO: validate
	comment.PostID = postID
	comment.Sender.ID = *payload.UserID
	comment.Content = *payload.Content

	err = app.store.Comments.Create(r.Context(), &comment)

	if errors.Is(err, store.ErrDataNotFound) {
		app.resourceNotFoundResponse(w, r, err)
		return
	} else if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if err := jsonDataResponse(w, http.StatusOK, comment); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}

func (app *application) getPostCommentsHandler(w http.ResponseWriter, r *http.Request) {
	postIDParam := chi.URLParam(r, "post_id")
	postID, err := strconv.ParseInt(postIDParam, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, ErrInvalidPostID)
		return
	}

	comments, err := app.store.Comments.GetByPostID(r.Context(), postID)

	if errors.Is(err, store.ErrDataNotFound) {
		app.resourceNotFoundResponse(w, r, err)
		return
	} else if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if err := jsonDataResponse(w, http.StatusOK, comments); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}

func (app *application) updatePostCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentIDParam := chi.URLParam(r, "comment_id")
	commentID, err := strconv.ParseInt(commentIDParam, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, ErrInvalidPostID)
		return
	}
	comment, err := app.store.Comments.GetByID(r.Context(), commentID)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	var payload CommentPayload
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if payload.Content != nil {
		comment.Content = *payload.Content
	}

	log.Println("comment", comment)

	if err := app.store.Comments.Update(r.Context(), comment); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if err := jsonDataResponse(w, http.StatusOK, comment); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}

func (app *application) deletePostCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentIDParam := chi.URLParam(r, "comment_id")
	commentID, err := strconv.ParseInt(commentIDParam, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, ErrInvalidPostID)
		return
	}

	err = app.store.Comments.DeleteByID(r.Context(), commentID)

	if errors.Is(err, store.ErrDataNotFound) {
		app.resourceNotFoundResponse(w, r, err)
		return
	} else if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
