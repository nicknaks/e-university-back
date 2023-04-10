package store

import (
	"back/graph/model"
	"context"
	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) ListFaculties(ctx context.Context) ([]*model.Faculty, error) {
	query, _, err := s.Builder().Select("*").From("faculties").ToSql()
	if err != nil {
		return nil, err
	}

	var res []*model.Faculty

	err = s.SelectContext(ctx, &res, query)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Storage) ListDepartments(ctx context.Context, facultiesID []string) ([]*model.Department, error) {
	query, args, err := s.Builder().Select("*").From("departments").Where(sq.Eq{"faculty_id": facultiesID}).ToSql()
	if err != nil {
		return nil, err
	}

	var res []*model.Department

	err = s.SelectContext(ctx, &res, query, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Storage) ListGroups(ctx context.Context, filter *model.GroupsFilter) ([]*model.Group, error) {
	query := s.Builder().Select("id, number, course").From("groups")

	if filter != nil {
		if filter.Course != nil {
			query = query.Where(sq.Eq{"course": *filter.Course})
		}
		if filter.DepartmentID != nil {
			query = query.Where(sq.Eq{"department_id": *filter.DepartmentID})
		}
	}

	queryStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var res []*model.Group

	err = s.SelectContext(ctx, &res, queryStr, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}
