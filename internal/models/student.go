package models

import (
	"back/graph/model"
)

type Student struct {
	ID      string
	GroupID string
	Name    string
}

func ToStudent(lesson *Student) *model.Student {
	if lesson == nil {
		return nil
	}

	return &model.Student{
		ID:      lesson.ID,
		Name:    lesson.Name,
		GroupID: lesson.GroupID,
	}
}
func ToStudents(lessons []*Student) []*model.Student {
	result := make([]*model.Student, 0, len(lessons))
	for _, lesson := range lessons {
		result = append(result, ToStudent(lesson))
	}
	return result
}
