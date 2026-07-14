package main

import (
	"net/http"

	"github.com/anikmahidul9/social/internal/store"
)

type CreateReplyPayload struct {
	Content string `json:"content" validate:"required,max=1000"`
}

func (app *application) createReplyHandler(w http.ResponseWriter, r *http.Request) {

	var payload CreateReplyPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(&payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	user := GetUserFromContext(r)
	parent := GetCommentFromCtx(r)

	reply := &store.Comment{
		PostID:          parent.PostID,
		UserID:          user.ID,
		ParentCommentID: &parent.ID,
		Content:         payload.Content,
	}

	err := app.store.Comments.CreateReplies(r.Context(), reply)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	reply.User = *user

	app.jsonResponse(w, http.StatusCreated, reply)
}

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required,max=1000"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {

	var payload CreateCommentPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(&payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	user := GetUserFromContext(r)
	post := getPostFromCtx(r)

	comment := &store.Comment{
		PostID:  post.ID,
		UserID:  user.ID,
		Content: payload.Content,
	}
	

	err := app.store.Comments.Create(r.Context(), comment)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	comment.User = *user

	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) likeCommentHandler(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)
	comment := GetCommentFromCtx(r)

	err := app.store.Reacts.Like(
		r.Context(),
		store.CommentLikeTable,
		store.CommentIDColumn,
		user.ID,
		comment.ID,
	)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) unlikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)
	comment := GetCommentFromCtx(r)

	err := app.store.Reacts.Unlike(
		r.Context(),
		store.CommentLikeTable,
		store.CommentIDColumn,
		user.ID,
		comment.ID,
	)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
