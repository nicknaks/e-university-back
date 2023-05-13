package store

import (
	"back/graph/model"
	"back/internal/models"
	"context"
	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) CreateSubject(ctx context.Context, input model.SubjectCreateInput) (*models.Subject, error) {
	query := s.Builder().Insert("subjects").SetMap(map[string]interface{}{
		"teacherid": input.TeacherID,
		"groupid":   input.GroupID,
		"name":      input.Name,
		"type":      models.ParseApiSubjectType(input.Type),
	}).Suffix("RETURNING *")

	sub := models.Subject{}

	err := s.Getx(ctx, &sub, query)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (s *Storage) UpdateSubject(ctx context.Context, subjectID string, m map[string]interface{}) (*models.Subject, error) {
	query := s.Builder().Update("subjects").SetMap(m).Where(sq.Eq{"id": subjectID}).Suffix("RETURNING *")

	sub := models.Subject{}

	err := s.Getx(ctx, &sub, query)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (s *Storage) LessonCreate(ctx context.Context, input model.LessonCreateInput, subject *models.Subject) (*models.Lesson, error) {
	query := s.Builder().Insert("lesson").SetMap(map[string]interface{}{
		"type":          models.ParseApiLessonType(input.Type),
		"subjectid":     input.SubjectID,
		"couple":        input.Couple,
		"day":           input.Day,
		"name":          subject.Name,
		"groupid":       subject.GroupID,
		"teacherid":     subject.TeacherID,
		"cabinet":       input.Cabinet,
		"isdenominator": input.IsDenominator,
		"isnumerator":   input.IsNumerator,
	}).Suffix("RETURNING *")

	lesson := models.Lesson{}

	err := s.Getx(ctx, &lesson, query)
	if err != nil {
		return nil, err
	}
	return &lesson, nil
}
