package store

import (
	"back/graph/model"
	"back/internal/auth_service"
	"back/internal/models"
	"context"
	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) GetToken(ctx context.Context, login string, pass string) (string, error) {
	query := s.Builder().Select("token").From("users").Where(sq.And{sq.Eq{"login": login}, sq.Eq{"password": pass}})

	var token string

	err := s.Getx(ctx, &token, query)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *Storage) GetUser(ctx context.Context, token string) (*auth_service.User, error) {
	query := s.Builder().Select("id, ownerId, type").From("users").Where(sq.Eq{"token": token})

	user := auth_service.User{}

	err := s.Getx(ctx, &user, query)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Storage) ListFaculties(ctx context.Context) ([]*model.Faculty, error) {
	query := s.Builder().Select("*").From("faculties")

	var res []*model.Faculty

	err := s.Selectx(ctx, &res, query)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Storage) ListDepartments(ctx context.Context, facultiesID []string) ([]*model.Department, error) {
	query := s.Builder().Select("*").From("departments")

	if len(facultiesID) != 0 {
		query = query.Where(sq.Eq{"facultyId": facultiesID})
	}

	var res []*model.Department

	err := s.Selectx(ctx, &res, query)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Storage) ListSubjects(ctx context.Context, filter *model.SubjectsFilter) ([]*models.Subject, error) {
	query := s.Builder().Select("*").From("subjects")

	if filter != nil {
		if filter.GroupID != nil {
			query = query.Where(sq.Eq{"groupid": *filter.GroupID})
		}
		if filter.TeacherID != nil {
			query = query.Where(sq.Eq{"teacherid": *filter.TeacherID})
		}
		if len(filter.ID) > 0 {
			query = query.Where(sq.Eq{"id": filter.ID})
		}
	}

	var res []*models.Subject

	err := s.Selectx(ctx, &res, query)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Storage) ListGroups(ctx context.Context, filter *model.GroupsFilter) ([]*model.Group, error) {
	query := s.Builder().Select("id, number, course").From("groups")

	if filter != nil {
		if len(filter.IDIn) > 0 {
			query = query.Where(sq.Eq{"id": filter.IDIn})
		}
		if filter.Course != nil {
			query = query.Where(sq.Eq{"course": *filter.Course})
		}
		if filter.DepartmentID != nil {
			query = query.Where(sq.Eq{"departmentId": *filter.DepartmentID})
		}
		if filter.IsMagistracy != nil {
			query = query.Where(sq.Eq{"ismagistracy": *filter.IsMagistracy})
		}
	}

	var res []*model.Group

	err := s.Selectx(ctx, &res, query)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Storage) ListLessons(ctx context.Context, filter model.ScheduleFilter) ([]*models.Lesson, error) {
	query := s.Builder().Select("*").From("lesson")

	if filter.GroupID != nil {
		query = query.Where(sq.Eq{"groupId": *filter.GroupID})
	}
	if filter.TeacherID != nil {
		query = query.Where(sq.Eq{"teacherId": *filter.TeacherID})
	}

	var res []*models.Lesson

	err := s.Selectx(ctx, &res, query)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Storage) ListTeachers(ctx context.Context, filter *model.TeachersFilter) ([]*models.Teacher, error) {
	query := s.Builder().Select("*").From("teachers")

	if filter != nil {
		if len(filter.IDIn) > 0 {
			query = query.Where(sq.Eq{"id": filter.IDIn})
		}
	}
	var res []*models.Teacher

	err := s.Selectx(ctx, &res, query)
	if err != nil {
		return nil, err
	}
	return res, nil
}
