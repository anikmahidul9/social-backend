package auth

import "context"

type contextKey string

const UserContextKey contextKey = "userID"

func SetUserID(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, UserContextKey, id)
}

func GetUserID(ctx context.Context) (int64, bool) {

	id, ok := ctx.Value(UserContextKey).(int64)

	return id, ok
}
