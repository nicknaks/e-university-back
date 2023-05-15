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

func (r *classResolver) StudentProgress(ctx context.Context, obj *model.Class) ([]*model.ClassProgress, error) {
	res, err := r.Storage.ListClassesProgresses(ctx, &model.ClassesProgressFilter{ClassID: lo.ToPtr(obj.ID)})
	if err != nil {
		return nil, err
	}

	return models.ToClassProgresses(res), nil
}

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

func (r *mutationResolver) SubjectCreate(ctx context.Context, input model.SubjectCreateInput) (*model.Subject, error) {
	subject, err := r.Storage.CreateSubject(ctx, input)
	if err != nil {
		return nil, err
	}

	return models.ToSubject(subject), nil
}

func (r *mutationResolver) SubjectTypeChange(ctx context.Context, input model.SubjectTypeChangeInput) (*model.Subject, error) {
	m := map[string]interface{}{
		"type": models.ParseApiSubjectType(input.Type),
	}

	subject, err := r.Storage.UpdateSubject(ctx, input.ID, m)
	if err != nil {
		return nil, err
	}

	return models.ToSubject(subject), nil
}

func (r *mutationResolver) LessonCreate(ctx context.Context, input model.LessonCreateInput) (*model.Lesson, error) {
	subjects, err := r.Storage.ListSubjects(ctx, &model.SubjectsFilter{ID: []string{input.SubjectID}})
	if err != nil {
		return nil, err
	}

	lesson, err := r.Storage.LessonCreate(ctx, input, subjects[0])
	return models.ToLesson(lesson), err
}

func (r *mutationResolver) StudentCreate(ctx context.Context, input model.StudentCreateInput) (*model.Student, error) {
	student, err := r.Storage.CreateStudent(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("Storage.CreateStudent err %w", err)
	}

	_, err = r.Storage.CreateClassesProgressForStudent(ctx, student.ID, student.GroupID)
	if err != nil {
		return nil, fmt.Errorf("Storage.CreateClassesProgressForStudent err %w", err)
	}

	_, err = r.Storage.CreateSubjectsResultsForStudent(ctx, student.ID, student.GroupID)
	if err != nil {
		return nil, fmt.Errorf("Storage.CreateSubjectsResultsForStudent err %w", err)
	}

	return models.ToStudent(student), nil
}

func (r *mutationResolver) MarkCreate(ctx context.Context, input model.MarkCreateInput) (*model.ClassProgress, error) {
	res, err := r.Storage.ClassProgressUpdate(ctx, input.ClassProgressID, map[string]interface{}{
		"mark": input.Mark,
	})
	if err != nil {
		return nil, err
	}

	return models.ToClassProgress(res), nil
}

func (r *mutationResolver) AbsentSet(ctx context.Context, input model.AbsentSetInput) ([]*model.ClassProgress, error) {
	res, err := r.Storage.SetAbsent(ctx, input.ClassProgressID)
	if err != nil {
		return nil, err
	}

	return models.ToClassProgresses(res), nil
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
	case 2:
		students, err := r.Storage.ListStudents(ctx, &model.StudentsFilter{
			IDIn: []string{user.OwnerID},
		})
		if err != nil {
			return nil, err
		}
		filter.GroupID = lo.ToPtr(students[0].GroupID)
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
	case 2:
		userType = model.UserTypeStudent
	case 3:
		userType = model.UserTypeAdmin
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

func (r *queryResolver) Students(ctx context.Context, filter *model.StudentsFilter) ([]*model.Student, error) {
	st, err := r.Storage.ListStudents(ctx, filter)
	if err != nil {
		return nil, err
	}

	return models.ToStudents(st), nil
}

func (r *queryResolver) Classes(ctx context.Context, filter *model.ClassesFilter) ([]*model.Class, error) {
	cl, err := r.Storage.ListClasses(ctx, filter)
	if err != nil {
		return nil, err
	}

	return models.ToClasses(cl), nil
}

func (r *queryResolver) SubjectResults(ctx context.Context, filter *model.SubjectResultsFilter) ([]*model.SubjectResult, error) {
	res, err := r.Storage.ListSubjectsResults(ctx, filter)
	if err != nil {
		return nil, err
	}

	return models.ToSubjectResults(res), nil
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

func (r *subjectResultResolver) Subject(ctx context.Context, obj *model.SubjectResult) ([]*model.Subject, error) {
	subjects, err := r.Storage.ListSubjects(ctx, &model.SubjectsFilter{ID: []string{obj.SubjectID}})
	if err != nil {
		return nil, err
	}

	return models.ToSubjects(subjects), nil
}

// Class returns generated.ClassResolver implementation.
func (r *Resolver) Class() generated.ClassResolver { return &classResolver{r} }

// Lesson returns generated.LessonResolver implementation.
func (r *Resolver) Lesson() generated.LessonResolver { return &lessonResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subject returns generated.SubjectResolver implementation.
func (r *Resolver) Subject() generated.SubjectResolver { return &subjectResolver{r} }

// SubjectResult returns generated.SubjectResultResolver implementation.
func (r *Resolver) SubjectResult() generated.SubjectResultResolver { return &subjectResultResolver{r} }

type classResolver struct{ *Resolver }
type lessonResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subjectResolver struct{ *Resolver }
type subjectResultResolver struct{ *Resolver }
