package token

import (
	"context"
	"net/http"
	"strings"
)

var ContextKey = "AuthToken"

func HTTPTokenToContext(ctx context.Context, r *http.Request) context.Context {
	token := r.Header.Get("Authorization")
	if strings.HasPrefix(token, "Bearer ") {
		token = token[7:]
	}
	return context.WithValue(ctx, ContextKey, token)
}

func HTTPTokenFromContext(ctx context.Context, r *http.Request) context.Context {
	token, ok := ctx.Value(ContextKey).(string)
	if ok {
		r.Header.Add("Authorization", "Bearer "+token)
	}
	return ctx
}
