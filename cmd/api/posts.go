package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/anikmahidul9/social/internal/store"
	"github.com/go-chi/chi/v5"
)

type Visibility string

const (
	Public  Visibility = "public"
	Private Visibility = "private"
)

type CreatePostPayload struct {
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	Visibility Visibility `json:"visibility"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseMultipartForm(20 << 20); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	user := GetUserFromContext(r)
	if user == nil {
		app.unauthorizedErrorResponse(w, r, errors.New("authentication required"))
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	visibility := store.Visibility(r.FormValue("visibility"))
	if visibility == "" {
		visibility = store.Public
	}

	post := &store.Post{
		Title:      title,
		Content:    content,
		Visibility: visibility,
		UserID:     user.ID,
	}

	if err := app.store.Posts.Create(r.Context(), post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	files := r.MultipartForm.File["images"]

	var imageURLs []string

	for _, file := range files {

		imageURL, err := saveImage(file)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		imageURLs = append(imageURLs, imageURL)
	}
	if len(imageURLs) > 0 {
		err := app.store.Images.Create(
			r.Context(),
			post.ID,
			imageURLs,
		)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
	}
}
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	user := GetUserFromContext(r)

	if post.Visibility == store.Private {
		if user == nil || user.ID != post.UserID {
			app.forbiddenResponse(w, r)
			return
		}
	}
	comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = comments

	likes, err := app.store.Reacts.GetLatestPostLikes(
		r.Context(),
		post.ID,
		5,
	)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.LatestLikes = likes

	count, err := app.store.Reacts.CountPostLikes(
		r.Context(),
		post.ID,
	)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.LikesCount = count
	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type UpdatePostPayload struct {
	Title      *string     `json:"title"`
	Content    *string     `json:"content"`
	Visibility *Visibility `json:"visibility"`
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	var payload UpdatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}

	if payload.Visibility != nil {
		post.Visibility = store.Visibility(*payload.Visibility)
	}

	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(postId, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	ctx := r.Context()
	if err := app.store.Posts.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}
	return
}

func (app *application) likePostHandler(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)
	post := getPostFromCtx(r)

	err := app.store.Reacts.Like(
		r.Context(),
		store.PostLikeTable,
		store.PostIDColumn,
		user.ID,
		post.ID,
	)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) unlikePostHandler(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)
	post := getPostFromCtx(r)

	err := app.store.Reacts.Unlike(
		r.Context(),
		store.PostLikeTable,
		store.PostIDColumn,
		user.ID,
		post.ID,
	)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postID")
		id, err := strconv.ParseInt(idParam, 10, 64)

		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()
		post, err := app.store.Posts.Get(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, "post", post)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *store.Post {
	post, _ := r.Context().Value("post").(*store.Post)

	return post
}
