package parser

import (
	"back/internal/store"
	"context"
	"github.com/dgryski/trifles/uuid"
	"github.com/sethvargo/go-password/password"
)

func InitUsers(ctx context.Context, db store.Storage) error {
	teachers, err := db.ListTeachers(ctx, nil)
	if err != nil {
		return err
	}

	for _, teacher := range teachers {
		query := db.Builder().Insert("users").SetMap(map[string]interface{}{
			"type":     1,
			"name":     teacher.Name,
			"login":    uuid.UUIDv4(),
			"password": password.MustGenerate(20, 10, 10, false, false),
			"ownerId":  teacher.ID,
			"token":    uuid.UUIDv4(),
		})
		err = db.Exec(ctx, query)
		if err != nil {
			return err
		}
	}

	return nil
}
