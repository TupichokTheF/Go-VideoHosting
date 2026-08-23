package app_context

import "context"

type contextKey int

const (
	_ contextKey = iota
	userIDKey
)

func ContextWithUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func UserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(userIDKey).(int)

	return userID, ok
}
