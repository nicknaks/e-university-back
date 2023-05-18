package models

import (
	"back/graph/model"
)

type SubjectResult struct {
	ID               string
	StudentID        string
	SubjectID        string
	FirstModuleMark  int
	SecondModuleMark int
	ThirdModuleMark  int
	Mark             int
	ExamResult       int
}

func ToSubjectResult(lesson *SubjectResult) *model.SubjectResult {
	if lesson == nil {
		return nil
	}

	return &model.SubjectResult{
		ID:               lesson.ID,
		StudentID:        lesson.StudentID,
		SubjectID:        lesson.SubjectID,
		FirstModuleMark:  lesson.FirstModuleMark,
		SecondModuleMark: lesson.SecondModuleMark,
		ThirdModuleMark:  lesson.ThirdModuleMark,
		Mark:             lesson.Mark,
		ExamResult:       lesson.ExamResult,
		CountAbsent:      2,
	}
}
func ToSubjectResults(lessons []*SubjectResult) []*model.SubjectResult {
	result := make([]*model.SubjectResult, 0, len(lessons))
	for _, lesson := range lessons {
		result = append(result, ToSubjectResult(lesson))
	}
	return result
}
