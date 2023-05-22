package parser

import (
	"back/internal/store"
	"context"
	"github.com/dgryski/trifles/uuid"
	"math/rand"
	"strings"
	"time"
)

func generatePass() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("QWERTYUIOPASDFGHJKLZXCVBNM" +
		"qwertyuiopasdfghjklzxcvbnm" +
		"0123456789")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

func generateLogin() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("QWERTYUIOPASDFGHJKLZXCVBNM" +
		"qwertyuiopasdfghjklzxcvbnm")
	length := 10
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

func InitUsers(ctx context.Context, db store.Storage) error {
	teachers, err := db.ListTeachers(ctx, nil)
	if err != nil {
		return err
	}

	for _, teacher := range teachers {
		query := db.Builder().Insert("users").SetMap(map[string]interface{}{
			"type":     1,
			"login":    generateLogin(),
			"password": generatePass(),
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
