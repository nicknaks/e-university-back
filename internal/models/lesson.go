package models

import "back/graph/model"

type Lesson struct {
	ID            string  `json:"id"`
	Type          int     `json:"type"`
	SubjectID     string  `json:"subjectID"`
	Name          *string `json:"name"`
	Couple        int     `json:"couple"`
	Day           int     `json:"day"`
	GroupID       string  `json:"groupID"`
	TeacherID     *string `json:"teacherID"`
	Cabinet       *string `json:"cabinet"`
	IsDenominator bool    `json:"isDenominator"`
	IsNumerator   bool    `json:"isNumerator"`
}

func ParseApiLessonType(lessonType model.LessonType) int {
	switch lessonType {
	case model.LessonTypeLab:
		return 1
	case model.LessonTypeLec:
		return 2
	case model.LessonTypeSem:
		return 3
	default:
		return 0
	}
}

func parseLessonType(lessonType int) model.LessonType {
	switch lessonType {
	case 0:
		return model.LessonTypeDefault
	case 1:
		return model.LessonTypeLab
	case 2:
		return model.LessonTypeLec
	case 3:
		return model.LessonTypeSem
	default:
		return model.LessonTypeDefault
	}
}

func ToLesson(lesson *Lesson) *model.Lesson {
	if lesson == nil {
		return nil
	}

	return &model.Lesson{
		ID:            lesson.ID,
		Type:          parseLessonType(lesson.Type),
		SubjectID:     lesson.SubjectID,
		Name:          lesson.Name,
		Couple:        lesson.Couple,
		Day:           lesson.Day,
		GroupID:       lesson.GroupID,
		TeacherID:     lesson.TeacherID,
		Cabinet:       lesson.Cabinet,
		IsDenominator: lesson.IsDenominator,
		IsNumerator:   lesson.IsNumerator,
	}
}

func ToLessons(lessons []*Lesson) []*model.Lesson {
	result := make([]*model.Lesson, 0, len(lessons))
	for _, lesson := range lessons {
		result = append(result, ToLesson(lesson))
	}
	return result
}
