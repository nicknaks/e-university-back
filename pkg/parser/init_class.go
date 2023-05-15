package parser

import (
	"back/graph/model"
	"back/internal/models"
	"back/internal/store"
	"context"
	"fmt"
	"gopkg.in/guregu/null.v4/zero"
	"time"
)

func InitParser(ctx context.Context, db store.Storage) error {
	lessons, err := db.ListLessons(ctx, model.ScheduleFilter{})
	if err != nil {
		return err
	}

	//понедельник
	lastDay := time.Date(2023, time.July, 6, 12, 0, 0, 0, time.UTC)

	for _, lesson := range lessons {
		firstMonday := time.Date(2023, time.February, 6, 12, 0, 0, 0, time.UTC) // знаменатель
		firstMonday = firstMonday.AddDate(0, 0, lesson.Day-1)
		for firstMonday.Before(lastDay) {
			// неделя - знаменатель
			if _, week := firstMonday.ISOWeek(); week%2 == 0 {
				// пара по числителям
				if lesson.IsNumerator == true {
					fmt.Println(firstMonday)
					firstMonday = firstMonday.AddDate(0, 0, 7)
					continue
				}
			} else {
				// пара по знаменателям
				if lesson.IsDenominator == true {
					fmt.Println(firstMonday)
					firstMonday = firstMonday.AddDate(0, 0, 7)
					continue
				}
			}

			class := models.Class{
				Day:       zero.TimeFrom(firstMonday),
				Type:      lesson.Type,
				Module:    getModule(firstMonday),
				SubjectID: lesson.SubjectID,
				LessonID:  lesson.ID,
				GroupID:   lesson.GroupID,
			}

			_, err = db.ClassCreate(ctx, class)
			if err != nil {
				return err
			}

			fmt.Println(firstMonday)
			firstMonday = firstMonday.AddDate(0, 0, 7)
		}
	}

	return nil
}

func getModule(date time.Time) int {
	if date.After(time.Date(2023, time.May, 1, 12, 0, 0, 0, time.UTC)) {
		return 3
	}

	if date.After(time.Date(2023, time.April, 1, 12, 0, 0, 0, time.UTC)) {
		return 2
	}

	return 1
}
