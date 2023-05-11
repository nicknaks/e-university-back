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

func (r *lessonResolver) Teacher(ctx context.Context, obj *model.Lesson) (*model.Teacher, error) {
	if obj.TeacherID == nil {
		return nil, nil
	}
	teachers, err := r.Storage.ListTeachers(ctx, &model.TeachersFilter{IDIn: []string{*obj.TeacherID}})
	if err != nil {
		return nil, err
	}
	return models.ToTeacher(teachers[0]), nil
}

func (r *lessonResolver) Group(ctx context.Context, obj *model.Lesson) (*model.Group, error) {
	teachers, err := r.Storage.ListGroups(ctx, &model.GroupsFilter{IDIn: []string{obj.GroupID}})
	if err != nil {
		return nil, err
	}

	return &model.Group{
		ID:     teachers[0].ID,
		Number: teachers[0].Number,
		Course: teachers[0].Course,
	}, nil
}

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

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	user := auth_service.GetUserFromContext(ctx)

	var userType model.UserType

	switch user.UserType {
	case 1:
		userType = model.UserTypeTeacher
	default:
		userType = model.UserTypeUnknown
	}

	return &model.User{
		ID:      user.ID,
		OwnerID: lo.ToPtr(user.OwnerID),
		Type:    userType,
	}, nil
}

func (r *queryResolver) Subjects(ctx context.Context, filter *model.SubjectsFilter) ([]*model.Subject, error) {
	subjects, err := r.Storage.ListSubjects(ctx, filter)
	if err != nil {
		fmt.Println("Storage.ListSubjects err %w", err)
	}

	return models.ToSubjects(subjects), nil
}

func (r *subjectResolver) Group(ctx context.Context, obj *model.Subject) (*model.Group, error) {
	teachers, err := r.Storage.ListGroups(ctx, &model.GroupsFilter{IDIn: []string{obj.GroupID}})
	if err != nil {
		return nil, err
	}

	return &model.Group{
		ID:     teachers[0].ID,
		Number: teachers[0].Number,
		Course: teachers[0].Course,
	}, nil
}

func (r *subjectResolver) Teacher(ctx context.Context, obj *model.Subject) (*model.Teacher, error) {
	if obj.TeacherID == nil {
		return nil, nil
	}
	teachers, err := r.Storage.ListTeachers(ctx, &model.TeachersFilter{IDIn: []string{*obj.TeacherID}})
	if err != nil {
		return nil, err
	}
	return models.ToTeacher(teachers[0]), nil
}

// Lesson returns generated.LessonResolver implementation.
func (r *Resolver) Lesson() generated.LessonResolver { return &lessonResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subject returns generated.SubjectResolver implementation.
func (r *Resolver) Subject() generated.SubjectResolver { return &subjectResolver{r} }

type lessonResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subjectResolver struct{ *Resolver }
