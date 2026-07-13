package main

import (
	"errors"
	"net/http"

	"github.com/anikmahidul9/social/internal/store"
)

type RegisterUserPayload struct {
	FirstName string `json:"first_name" validate:"required,max=100"`
	LastName  string `json:"last_name" validate:"required,max=100"`
	Username  string `json:"username" validate:"required,max=100"`
	Email     string `json:"email" validate:"required,email,max=255"`
	Password  string `json:"password" validate:"required,min=4,max=72"`
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	user := &store.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Username:  payload.Username,
		Email:     payload.Email,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.store.Users.Create(r.Context(), user); err != nil {
		switch err {
		case store.ErrDuplicateUsername:
			app.conflictResponse(w, r, "username already exists")
			return

		case store.ErrDuplicateEmail:
			app.conflictResponse(w, r, "email already exists")
			return

		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusCreated, user); err != nil {
		app.internalServerError(w, r, err)
	}

}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {

	var payload LoginUserPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	user, err := app.store.Users.GetByEmail(r.Context(), payload.Email)

	if err != nil {

		if errors.Is(err, store.ErrNotFound) {
			app.unauthorizedErrorResponse(w, r, errors.New("invalid email or password"))
			return
		}

		app.internalServerError(w, r, err)
		return
	}

	err = user.Password.Matches(payload.Password)

	if err != nil {
		app.unauthorizedErrorResponse(w, r, errors.New("invalid email or password"))
		return
	}

	token, err := app.jwt.GenerateToken(user.ID)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(
		w,
		http.StatusOK,
		map[string]string{
			"token": token,
		},
	)
}
