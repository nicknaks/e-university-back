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
	Type      int
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
		Type:      parseSubjectType(lesson.Type),
	}
}

func ParseApiSubjectType(lessonType model.SubjectType) int {
	switch lessonType {
	case model.SubjectTypeUnknown:
		return 0
	case model.SubjectTypeCredit:
		return 1
	case model.SubjectTypeExam:
		return 2
	case model.SubjectTypeCourseWork:
		return 3
	case model.SubjectTypePractical:
		return 4
	default:
		return 0
	}
}

func parseSubjectType(lessonType int) model.SubjectType {
	switch lessonType {
	case 1:
		return model.SubjectTypeCredit
	case 2:
		return model.SubjectTypeExam
	case 3:
		return model.SubjectTypeCourseWork
	case 4:
		return model.SubjectTypePractical
	default:
		return model.SubjectTypeUnknown
	}
}

func ToSubjects(lessons []*Subject) []*model.Subject {
	result := make([]*model.Subject, 0, len(lessons))
	for _, lesson := range lessons {
		result = append(result, ToSubject(lesson))
	}
	return result
}
