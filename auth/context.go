package auth

import "context"

type ctxKey string

const userIDKey ctxKey = "userID"

// context 用ヘルパー
func WithUserID(ctx context.Context, id int64) context.Context {
    return context.WithValue(ctx, userIDKey, id)
}

func GetUserID(ctx context.Context) (int64, bool) {
    id, ok := ctx.Value(userIDKey).(int64)
    return id, ok
}
