package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Reensef/golang-social/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr string
	db   dbConfig
	env  string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

var (
	ErrInvalidCommentID = errors.New("invalid comment ID")
	ErrInvalidPostID    = errors.New("invalid post ID")
	ErrInvalidLimit     = errors.New("invalid limit")
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)

			r.Route("/{post_id}", func(r chi.Router) {
				r.Use(app.postContextMiddleware) // Need remove
				r.Get("/", app.getPostHandler)
				r.Delete("/", app.deletePostHandler)
				r.Patch("/", app.updatePostHandler)
			})

			r.Route("/{post_id}/comments", func(r chi.Router) {
				r.Post("/", app.createPostCommentHandler)
				r.Get("/", app.getPostCommentsHandler)
				r.Patch("/{comment_id}", app.updatePostCommentHandler)
				r.Delete("/{comment_id}", app.deletePostCommentHandler)
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Post("/", app.createUserHandler)
			// r.Get("/", app.getUsersHandler)
		// 	r.Patch("/{userID}", app.updateUserHandler)
		// 	r.Delete("/{userID}", app.deleteUserHandler)
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	srv := http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server has started as %s", app.config.addr)

	return srv.ListenAndServe()
}