package models

import (
	"back/graph/model"
	"gopkg.in/guregu/null.v4/zero"
)

type SubjectResult struct {
	ID                      string
	StudentID               string
	SubjectID               string
	FirstModuleMark         float64
	SecondModuleMark        float64
	ThirdModuleMark         float64
	Mark                    float64
	ExamResult              int
	FirstModuleMarkComment  zero.String
	SecondModuleMarkComment zero.String
	ThirdModuleMarkComment  zero.String
	ExamResultComment       zero.String
}

func ToSubjectResult(lesson *SubjectResult) *model.SubjectResult {
	if lesson == nil {
		return nil
	}

	return &model.SubjectResult{
		ID:                      lesson.ID,
		StudentID:               lesson.StudentID,
		SubjectID:               lesson.SubjectID,
		FirstModuleMark:         lesson.FirstModuleMark,
		SecondModuleMark:        lesson.SecondModuleMark,
		ThirdModuleMark:         lesson.ThirdModuleMark,
		Mark:                    lesson.Mark,
		ExamResult:              lesson.ExamResult,
		CountAbsent:             2,
		FirstModuleMarkComment:  lesson.FirstModuleMarkComment.Ptr(),
		SecondModuleMarkComment: lesson.SecondModuleMarkComment.Ptr(),
		ThirdModuleMarkComment:  lesson.ThirdModuleMarkComment.Ptr(),
		ExamResultComment:       lesson.ExamResultComment.Ptr(),
	}
}
func ToSubjectResults(lessons []*SubjectResult) []*model.SubjectResult {
	result := make([]*model.SubjectResult, 0, len(lessons))
	for _, lesson := range lessons {
		result = append(result, ToSubjectResult(lesson))
	}
	return result
}
