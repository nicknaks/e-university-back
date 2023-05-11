package models

import (
	"back/graph/model"
	"gopkg.in/guregu/null.v4/zero"
)

type Teacher struct {
	ID         string
	Name       zero.String
	Speciality zero.String
}

func ToTeacher(lesson *Teacher) *model.Teacher {
	if lesson == nil {
		return nil
	}

	return &model.Teacher{
		ID:   lesson.ID,
		Name: lesson.Name.Ptr(),
	}
}

func ToTeachers(lessons []*Teacher) []*model.Teacher {
	result := make([]*model.Teacher, 0, len(lessons))
	for _, lesson := range lessons {
		result = append(result, ToTeacher(lesson))
	}
	return result
}
