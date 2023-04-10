package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"back/graph/generated"
	"back/graph/model"
	"context"
	"fmt"
)

func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateTodov(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	panic(fmt.Errorf("not implemented"))
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
func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}
