package main

import (
	"context"
	"net/http"
)

type contextKey string

const userContextKey = contextKey("user_id")

func (app *application) contextSetUser(r *http.Request, userID int64) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, userID)
	return r.WithContext(ctx)
}

func (app *application) ContextGetUser(r *http.Request) (int64, bool) {
	userID, ok := r.Context().Value(userContextKey).(int64)
	return userID, ok
}