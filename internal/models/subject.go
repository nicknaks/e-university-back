package models

import (
	"back/graph/model"
	"gopkg.in/guregu/null.v4/zero"
)

type Subject struct {
	ID        string
	TeacherID zero.String
	GroupID   string
	Name      zero.String
}

func ToSubject(lesson *Subject) *model.Subject {
	if lesson == nil {
		return nil
	}

	return &model.Subject{
		ID:        lesson.ID,
		TeacherID: lesson.TeacherID.Ptr(),
		GroupID:   lesson.GroupID,
		Name:      lesson.Name.Ptr(),
	}
}

func ToSubjects(lessons []*Subject) []*model.Subject {
	result := make([]*model.Subject, 0, len(lessons))
	for _, lesson := range lessons {
		result = append(result, ToSubject(lesson))
	}
	return result
}
