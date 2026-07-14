package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/anikmahidul9/social/internal/store"
	"github.com/go-chi/chi/v5"
)

type contextKey string

const userCtx contextKey = "user"

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authorization := r.Header.Get("Authorization")

		if authorization == "" {
			app.unauthorizedErrorResponse(w, r, errors.New("missing authorization header"))
			return
		}

		parts := strings.Split(authorization, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorizedErrorResponse(w, r, errors.New("invalid authorization header"))
			return
		}

		token := parts[1]

		claims, err := app.jwt.ValidateToken(token)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		user, err := app.store.Users.GetById(r.Context(), claims.UserID)
		if err != nil {

			if errors.Is(err, store.ErrNotFound) {
				app.unauthorizedErrorResponse(w, r, errors.New("user not found"))
				return
			}

			app.internalServerError(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), userCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(r *http.Request) *store.User {

	user, ok := r.Context().Value(userCtx).(*store.User)

	if !ok {
		return nil
	}

	return user
}

func (app *application) postOwnerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := GetUserFromContext(r)
		if user == nil {
			app.unauthorizedErrorResponse(w, r, errors.New("authentication required"))
			return
		}

		post := getPostFromCtx(r)
		if post.UserID != user.ID {
			app.forbiddenResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})

}

const (
	userCtxKey    contextKey = "user"
	postCtxKey    contextKey = "post"
	commentCtxKey contextKey = "comment"
)

func GetCommentFromCtx(r *http.Request) *store.Comment {

	comment, ok := r.Context().Value(commentCtxKey).(*store.Comment)

	if !ok {
		panic("missing comment context")
	}

	return comment
}
func (app *application) commentContextMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id, err := strconv.ParseInt(
			chi.URLParam(r, "commentID"),
			10,
			64,
		)

		if err != nil {
			app.badRequestError(w, r, err)
			return
		}

		comment, err := app.store.Comments.GetByID(r.Context(), id)
		if err != nil {
			app.notFoundResponse(w, r, err)
			return
		}

		ctx := context.WithValue(
			r.Context(),
			commentCtxKey,
			comment,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
