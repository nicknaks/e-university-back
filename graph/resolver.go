package graph

import (
	"back/internal/store"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Storage store.Storage
}

func NewResolver() *Resolver {
	return &Resolver{
		Storage: store.NewStorage(),
	}
}
