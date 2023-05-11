package auth_service

import "context"

type UserKey string

type User struct {
	ID       string
	OwnerID  string
	UserType int `db:"type"`
}

const UserKeyContext = HTTPKey("USER")

func SetUserToContext(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, UserKeyContext, user)
}

func GetUserFromContext(ctx context.Context) User {
	val := ctx.Value(UserKeyContext)
	if val == nil {
		panic("http context not set")
	}
	httpCtx := ctx.Value(UserKeyContext).(User)
	return httpCtx
}
