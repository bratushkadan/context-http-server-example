package internal

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bratushkadan/context-http-server-example/pkg/constants"
	"github.com/google/uuid"
)

type Middleware = func(http.Handler) http.Handler

type authKey struct{}

var key authKey

func ContextWithAuth(ctx context.Context, authHeader string) context.Context {
	return context.WithValue(ctx, key, authHeader)
}

func AuthFromContext(ctx context.Context) (string, bool) {
	authHeader, ok := ctx.Value(key).(string)
	return authHeader, ok
}

func CreateAuthMiddleware() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authToken, ok := r.Header[constants.XAuthToken]
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(fmt.Sprintf("no %s provided : authenticated", constants.XAuthToken)))
				return
			}
			authTokenStr := strings.Join(authToken, ",")
			ctx := r.Context()
			_, ok = LookupUserNameByAuthToken(ctx, authTokenStr)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(fmt.Sprintf("invalid %s provided : authenticated", constants.XAuthToken)))
				return
			}
			ctx = ContextWithAuth(ctx, authTokenStr)
			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
	}
}

type guidKey int

const gKey guidKey = 1

func contextWithGUID(ctx context.Context, guid string) context.Context {
	return context.WithValue(ctx, gKey, guid)
}
func guidFromContext(ctx context.Context) (string, bool) {
	g, ok := ctx.Value(key).(string)
	return g, ok
}

// constants.XRequestId header is propagated in http requests to other services in http client
func CreateRequestIdMiddleware() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			ctx := req.Context()
			if guid := req.Header.Get(constants.XRequestId); guid != "" {
				ctx = contextWithGUID(ctx, guid)
			} else {
				ctx = contextWithGUID(ctx, uuid.New().String())
			}
			req = req.WithContext(ctx)
			h.ServeHTTP(rw, req)
		})
	}
}

func CreateTimeoutMiddleware(timeout time.Duration) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
	}
}
