package main

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/anikmahidul9/social/internal/store"
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
