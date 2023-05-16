package store

import (
	"back/graph/model"
	"back/internal/models"
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/dgryski/trifles/uuid"
	"github.com/samber/lo"
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

func (s *Storage) UpdateSubjectResult(ctx context.Context, studentID string, subjectID string, m map[string]interface{}) (*models.SubjectResult, error) {
	query := s.Builder().Update("subjects_results").SetMap(m).
		Where(sq.And{sq.Eq{"studentid": studentID}, sq.Eq{"subjectid": subjectID}}).
		Suffix("RETURNING *")

	sub := models.SubjectResult{}

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

func (s *Storage) ClassCreate(ctx context.Context, class models.Class) (*models.Class, error) {
	query := s.Builder().Insert("classes").SetMap(map[string]interface{}{
		"day":       class.Day,
		"type":      class.Type,
		"comment":   class.Comment,
		"name":      class.Name,
		"subtype":   class.SubType,
		"module":    class.Module,
		"subjectId": class.SubjectID,
		"lessonId":  class.LessonID,
		"groupId":   class.GroupID,
	}).Suffix("RETURNING *")

	lesson := models.Class{}

	err := s.Getx(ctx, &lesson, query)
	if err != nil {
		return nil, err
	}
	return &lesson, nil
}

func (s *Storage) CreateClassesProgressForStudent(ctx context.Context, studentID string, groupID string) ([]*models.ClassProgress, error) {
	classes, err := s.ListClasses(ctx, &model.ClassesFilter{GroupID: lo.ToPtr(groupID)})
	if err != nil {
		return nil, fmt.Errorf("ListClasses err %w", err)
	}

	for _, class := range classes {
		_, err = s.ClassProgressCreate(ctx, class.ID, studentID)
		if err != nil {
			return nil, fmt.Errorf("ClassProgressCreate err %w", err)
		}
	}

	return nil, nil
}

func (s *Storage) CreateSubjectsResultsForStudent(ctx context.Context, studentID string, groupID string) ([]*models.SubjectResult, error) {
	subjects, err := s.ListSubjects(ctx, &model.SubjectsFilter{GroupID: lo.ToPtr(groupID)})
	if err != nil {
		return nil, fmt.Errorf("ListSubjects err %w", err)
	}

	for _, subject := range subjects {
		_, err = s.SubjectResultCreate(ctx, subject.ID, studentID)
		if err != nil {
			return nil, fmt.Errorf("SubjectResultCreate err %w", err)
		}
	}

	return nil, nil
}

func (s *Storage) SubjectResultCreate(ctx context.Context, subjectID string, studentID string) (*models.SubjectResult, error) {
	query := s.Builder().Insert("subjects_results").SetMap(map[string]interface{}{
		"subjectid": subjectID,
		"studentid": studentID,
	}).Suffix("RETURNING *")

	lesson := models.SubjectResult{}

	err := s.Getx(ctx, &lesson, query)
	if err != nil {
		return nil, err
	}
	return &lesson, nil
}

func (s *Storage) ClassProgressCreate(ctx context.Context, classID string, studentID string) (*models.ClassProgress, error) {
	query := s.Builder().Insert("classes_progresses").SetMap(map[string]interface{}{
		"classid":   classID,
		"studentid": studentID,
	}).Suffix("RETURNING *")

	lesson := models.ClassProgress{}

	err := s.Getx(ctx, &lesson, query)
	if err != nil {
		return nil, err
	}
	return &lesson, nil
}

func (s *Storage) ClassProgressUpdate(ctx context.Context, classProgressID string, m map[string]interface{}) (*models.ClassProgress, error) {
	query := s.Builder().Update("classes_progresses").SetMap(m).Where(sq.Eq{"id": classProgressID}).Suffix("RETURNING *")

	lesson := models.ClassProgress{}

	err := s.Getx(ctx, &lesson, query)
	if err != nil {
		return nil, err
	}
	return &lesson, nil
}

func (s *Storage) SetAbsent(ctx context.Context, classProgressIDs []string) ([]*models.ClassProgress, error) {
	query := s.Builder().Update("classes_progresses").SetMap(map[string]interface{}{
		"isabsent": true,
	}).Where(sq.Eq{"id": classProgressIDs}).Suffix("RETURNING *")

	var res []*models.ClassProgress

	err := s.Selectx(ctx, &res, query)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Storage) SetAttended(ctx context.Context, classProgressIDs []string) ([]*models.ClassProgress, error) {
	query := s.Builder().Update("classes_progresses").SetMap(map[string]interface{}{
		"isabsent": false,
	}).Where(sq.Eq{"id": classProgressIDs}).Suffix("RETURNING *")

	var res []*models.ClassProgress

	err := s.Selectx(ctx, &res, query)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Storage) CreateStudent(ctx context.Context, input model.StudentCreateInput) (*models.Student, error) {
	query := s.Builder().Insert("students").SetMap(map[string]interface{}{
		"groupid": input.GroupID,
		"name":    input.Name,
	}).Suffix("RETURNING *")

	sub := models.Student{}

	err := s.Getx(ctx, &sub, query)
	if err != nil {
		return nil, err
	}

	query = s.Builder().Insert("users").SetMap(map[string]interface{}{
		"type":     2,
		"login":    uuid.UUIDv4(),
		"password": "123",
		"ownerid":  sub.ID,
		"token":    uuid.UUIDv4(),
	})

	err = s.Exec(ctx, query)
	if err != nil {
		return nil, err
	}

	return &sub, nil
}
