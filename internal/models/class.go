package models

import (
	"back/graph/model"
	"gopkg.in/guregu/null.v4/zero"
)

type Class struct {
	ID        string
	Day       zero.Time
	Type      int
	Comment   zero.String
	Name      zero.String
	SubType   int
	Module    int
	SubjectID string
	LessonID  string
	GroupID   string
}

func ToClass(lesson *Class) *model.Class {
	if lesson == nil {
		return nil
	}

	return &model.Class{
		ID:        lesson.ID,
		Day:       lesson.Day.Time.Format("02.01"),
		Type:      parseLessonType(lesson.Type),
		Comment:   lesson.Comment.Ptr(),
		Name:      lesson.Name.Ptr(),
		Module:    lesson.Module,
		SubjectID: lesson.SubjectID,
		LessonID:  lesson.LessonID,
		GroupID:   lesson.GroupID,
	}
}
func ToClasses(lessons []*Class) []*model.Class {
	result := make([]*model.Class, 0, len(lessons))
	for _, lesson := range lessons {
		result = append(result, ToClass(lesson))
	}
	return result
}
