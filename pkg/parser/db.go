package parser

import (
	"back/internal/store"
	"context"
	"fmt"
	"strings"
)

func InitData(ctx context.Context, db store.Storage, data []Fac) ([]Fac, error) {
	// insert faculties
	for i, f := range data {
		query := db.Builder().Insert("faculties").SetMap(map[string]interface{}{
			"number": f.Number,
			"name":   f.Name,
		}).Suffix("RETURNING id")

		var id string
		err := db.Getx(ctx, &id, query)
		if err != nil {
			return nil, err
		}

		for j, dep := range f.Deps {
			query = db.Builder().Insert("departments").SetMap(map[string]interface{}{
				"number":    dep.Number,
				"name":      " ",
				"facultyId": id,
			}).Suffix("RETURNING id")

			var depId string
			err = db.Getx(ctx, &depId, query)
			if err != nil {
				return nil, err
			}

			for k, group := range dep.Groups {
				query = db.Builder().Insert("groups").SetMap(map[string]interface{}{
					"number":       group.Number,
					"course":       getCourse(group.Number),
					"isMagistracy": getMagistracy(group.Number),
					"departmentId": depId,
				}).Suffix("RETURNING id")

				var groupID string
				err = db.Getx(ctx, &groupID, query)
				data[i].Deps[j].Groups[k].ID = groupID
				if err != nil {
					return nil, err
				}
			}
		}
	}

	fmt.Println("data inserted")

	return data, nil
}

func getCourse(number string) int {
	number = strings.SplitAfter(number, "-")[1]

	switch {
	case strings.HasPrefix(number, "2"):
		return 1
	case strings.HasPrefix(number, "4"):
		return 2
	case strings.HasPrefix(number, "6"):
		return 3
	case strings.HasPrefix(number, "8"):
		return 4
	case strings.HasPrefix(number, "10"):
		return 5
	default:
		panic(number)
	}

	return 0
}

func getMagistracy(number string) bool {
	return strings.HasSuffix(number, "М")
}
