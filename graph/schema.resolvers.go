package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"back/graph/generated"
	"back/graph/model"
	"back/internal/auth_service"
	"back/internal/models"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/samber/lo"
)

func (r *mutationResolver) Login(ctx context.Context, login string, password string) (bool, error) {
	token, err := r.Storage.GetToken(ctx, login, password)
	if err != nil {
		return false, err
	}

	httpContext := auth_service.GetHttpContext(ctx)

	http.SetCookie(httpContext.W, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(time.Hour * 24),
	})

	// ставим куки
	return true, nil
}

func (r *mutationResolver) Logout(ctx context.Context) (bool, error) {
	httpContext := auth_service.GetHttpContext(ctx)

	http.SetCookie(httpContext.W, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now().Add(time.Hour * (-1)),
	})

	return true, nil
}

func (r *queryResolver) Faculties(ctx context.Context) ([]*model.Faculty, error) {
	res, err := r.Storage.ListFaculties(ctx)
	if err != nil {
		return nil, fmt.Errorf("Storage.ListFaculties(ctx) err %w", err)
	}
	var ids []string

	for _, re := range res {
		ids = append(ids, re.ID)
	}

	deps, err := r.Storage.ListDepartments(ctx, ids)
	if err != nil {
		return nil, err
	}

	for _, dep := range deps {
		for i, re := range res {
			if re.ID == dep.FacultyID {
				res[i].Departments = append(res[i].Departments, dep)
			}
		}
	}

	return res, nil
}

func (r *queryResolver) Groups(ctx context.Context, filter *model.GroupsFilter) ([]*model.Group, error) {
	return r.Storage.ListGroups(ctx, filter)
}

func (r *queryResolver) Schedule(ctx context.Context, filter model.ScheduleFilter) ([]*model.Lesson, error) {
	lessons, err := r.Storage.ListLessons(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("Storage.ListLessons err %w", err)
	}

	return models.ToLessons(lessons), nil
}

func (r *queryResolver) Teachers(ctx context.Context, filter *model.TeachersFilter) ([]*model.Teacher, error) {
	lessons, err := r.Storage.ListTeachers(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("Storage.ListLessons err %w", err)
	}

	return models.ToTeachers(lessons), nil
}

func (r *queryResolver) MySchedule(ctx context.Context) ([]*model.Lesson, error) {
	user := auth_service.GetUserFromContext(ctx)

	filter := model.ScheduleFilter{}
	switch user.UserType {
	case 1:
		filter.TeacherID = lo.ToPtr(user.OwnerID)
	default:
		panic(nil)
	}

	lessons, err := r.Storage.ListLessons(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("Storage.ListLessons err %w", err)
	}

	return models.ToLessons(lessons), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *queryResolver) T(ctx context.Context) ([]*bool, error) {
	panic(fmt.Errorf("not implemented"))
}
