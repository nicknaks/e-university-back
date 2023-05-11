package parser

import (
	"back/internal/store"
	"context"
	"fmt"
	"strings"
)

func ExtractTeachers(ctx context.Context, db store.Storage, lessons []Lesson) error {
	teachers := map[string]string{}
	// store teachers
	for i, lesson := range lessons {
		var id string
		var ok bool

		switch lesson.Name {
		case "Самостоятельная работа":
			continue
		}

		if strings.ContainsRune(lesson.Teacher, ',') {
			lesson.Teacher = strings.TrimSpace(strings.TrimSuffix(strings.SplitAfter(lesson.Teacher, ",")[0], ","))
		}

		if id, ok = teachers[lesson.Teacher]; !ok {
			query := db.Builder().Insert("teachers").SetMap(map[string]interface{}{
				"name": lesson.Teacher,
			}).Suffix("Returning id")

			err := db.Getx(ctx, &id, query)

			if err != nil {
				return err
			}
			teachers[lesson.Teacher] = id
		}
		lessons[i].TeacherID = id
	}

	// store subjects
	for i, lesson := range lessons {
		var id string

		switch lesson.Name {
		case "Самостоятельная работа":
			continue
		}

		query := db.Builder().Insert("subjects").SetMap(map[string]interface{}{
			"teacherId": lesson.TeacherID,
			"groupId":   lesson.GroupID,
			"name":      lesson.Name,
		}).Suffix("ON CONFLICT (groupId, name)").Suffix("DO UPDATE SET name=EXCLUDED.name RETURNING id")

		err := db.Getx(ctx, &id, query)

		if err != nil {
			return err
		}

		lessons[i].SubjectID = id
	}

	// store lessons
	for _, lesson := range lessons {
		switch lesson.Name {
		case "Самостоятельная работа":
			continue
		}

		query := db.Builder().Insert("lesson").SetMap(map[string]interface{}{
			"type":          parseType(lesson.Type),
			"subjectId":     lesson.SubjectID,
			"couple":        parseCouple(lesson.Couple),
			"day":           parseDay(lesson.Day),
			"groupId":       lesson.GroupID,
			"teacherId":     lesson.TeacherID,
			"cabinet":       lesson.Place,
			"isDenominator": lesson.IsDenominator,
			"isNumerator":   lesson.IsNumerator,
			"name":          lesson.Name,
		})
		err := db.Exec(ctx, query)
		if err != nil {
			fmt.Println(lesson)
			return err
		}
	}

	return nil
}

func parseCouple(couple string) int {
	switch strings.TrimSpace(couple) {
	case "08:30 - 10:05":
		return 1
	case "10:15 - 11:50":
		return 2
	case "12:00 - 13:35":
		return 3
	case "13:50 - 15:25":
		return 4
	case "15:40 - 17:15":
		return 5
	case "17:25 - 19:00":
		return 6
	case "19:10 - 20:45":
		return 7
	default:
		panic(couple)
	}
}

func parseDay(day string) int {
	switch strings.TrimSpace(day) {
	case "ПН":
		return 1
	case "ВТ":
		return 2
	case "СР":
		return 3
	case "ЧТ":
		return 4
	case "ПТ":
		return 5
	case "СБ":
		return 6
	default:
		panic(day)
	}
}

func parseType(typeLesson string) int {
	switch strings.TrimSpace(typeLesson) {
	case "":
		return 0
	case "(лаб)":
		return 1
	case "(лек)":
		return 2
	case "(сем)":
		return 3
	default:
		panic(typeLesson)
	}
}
