package models

import (
	"back/graph/model"
	"gopkg.in/guregu/null.v4/zero"
)

type ClassProgress struct {
	ID        string
	ClassID   string
	StudentID string
	IsAbsent  bool
	TeacherID zero.String
	Mark      int
}

func ToClassProgress(lesson *ClassProgress) *model.ClassProgress {
	if lesson == nil {
		return nil
	}

	return &model.ClassProgress{
		ID:        lesson.ID,
		ClassID:   lesson.ClassID,
		StudentID: lesson.StudentID,
		IsAbsent:  lesson.IsAbsent,
		TeacherID: lesson.TeacherID.Ptr(),
		Mark:      lesson.Mark,
	}
}
func ToClassProgresses(lessons []*ClassProgress) []*model.ClassProgress {
	result := make([]*model.ClassProgress, 0, len(lessons))
	for _, lesson := range lessons {
		result = append(result, ToClassProgress(lesson))
	}
	return result
}
