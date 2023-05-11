package store

import (
	"back/internal/auth_service"
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
)

func (db *Storage) IsAuth(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	httpContext := auth_service.GetHttpContext(ctx)
	cookie, err := httpContext.R.Cookie("token")
	if err != nil {
		return nil, fmt.Errorf("cookie not setted")
	}

	user, err := db.GetUser(ctx, cookie.Value)
	if err != nil {
		return nil, err
	}

	ctx = auth_service.SetUserToContext(ctx, *user)

	return next(ctx)
}
