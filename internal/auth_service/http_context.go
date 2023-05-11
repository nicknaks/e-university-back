package auth_service

import (
	"context"
	"net/http"
)

type HTTPKey string

type HTTP struct {
	W http.ResponseWriter
	R *http.Request
}

const httpKeyContext = HTTPKey("HTTP")

func GetHttpContext(ctx context.Context) *HTTP {
	val := ctx.Value(httpKeyContext)
	if val == nil {
		panic("http context not set")
	}
	httpCtx := ctx.Value(httpKeyContext).(HTTP)
	return &httpCtx
}

// WithHttpContext set HTTP on context
func WithHttpContext(ctx context.Context, h HTTP) context.Context {
	return context.WithValue(ctx, httpKeyContext, h)
}

// InjectHTTPMiddleware handles injecting the ResponseWriter and Request structs
// into context so that resolver methods can use these to set and read cookies.
func InjectHTTPMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			httpContext := HTTP{
				W: w,
				R: r,
			}
			ctx := context.WithValue(r.Context(), httpKeyContext, httpContext)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
