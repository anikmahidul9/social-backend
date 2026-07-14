package main

import (
	"log"
	"net/http"
	"time"

	"github.com/anikmahidul9/social/internal/auth"
	"github.com/anikmahidul9/social/internal/store"
	"github.com/go-chi/chi/v5"

	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
	store  store.Storage
	jwt    *auth.JWTAuthenticator
}

type config struct {
	addr string
	db   dbConfig
	auth authConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}
type authConfig struct {
	secret string
}

func (app *application) mount() *chi.Mux {
	r := chi.NewRouter()
	r.Use(app.corsMiddleware)
	r.Use(middleware.Logger)
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		r.Route("/posts", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/", app.createPostHandler)
			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)

				// No owner check
				r.Put("/like", app.likePostHandler)
				r.Delete("/like", app.unlikePostHandler)
				r.Get("/", app.getPostHandler)
				r.Post("/comments", app.createCommentHandler)

				// Owner-only routes
				r.Group(func(r chi.Router) {
					r.Use(app.postOwnerMiddleware)

					r.Patch("/", app.updatePostHandler)
					r.Delete("/", app.deletePostHandler)
				})
			})

		})
		r.Route("/users", func(r chi.Router) {
			//	r.Post("/", app.createPostHandler)
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.usersContextMiddleware)
				r.Get("/", app.getUserHandler)
				// r.Patch("/", app.updatePostHandler)
				// r.Delete("/", app.deletePostHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/feed", app.getUserFeedHandler)
			})
		})
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
			r.Post("/login", app.loginHandler)
		})

		r.Route("/comments", func(r chi.Router) {

			r.Use(app.AuthTokenMiddleware)

			r.Route("/{commentID}", func(r chi.Router) {

				r.Use(app.commentContextMiddleware)

				r.Post("/replies", app.createReplyHandler)

				// r.Patch("/", app.updateCommentHandler)
				// r.Delete("/", app.deleteCommentHandler)

				r.Put("/like", app.likeCommentHandler)
				r.Delete("/like", app.unlikeCommentHandler)
			})
		})

	})
	fs := http.FileServer(http.Dir("./uploads"))

	r.Handle("/uploads/*", http.StripPrefix("/uploads/", fs))
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
