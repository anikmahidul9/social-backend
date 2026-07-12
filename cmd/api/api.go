package main

import (
	"log"
	"net/http"
	"time"

	"github.com/anikmahidul9/social/internal/store"
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
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)
			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)
				r.Get("/", app.getPostHandler)
				r.Patch("/", app.updatePostHandler)
				r.Delete("/", app.deletePostHandler)
			})
		})
		r.Route("/users", func(r chi.Router) {
		//	r.Post("/", app.createPostHandler)
			r.Route("/{userID}", func(r chi.Router) {
				r.Get("/", app.getUserHandler)
				// r.Patch("/", app.updatePostHandler)
				// r.Delete("/", app.deletePostHandler)
			})
		})
	})
	return r
}
func (app *application) run(mux *chi.Mux) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server has started at %s", app.config.addr)
	return srv.ListenAndServe()

}
