package graph

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"

	"github.com/ZhanLiangUF/graphql-set/pg"
)

type Resolver struct {
	Repository pg.Repository
}

func (r *mutationResolver) CreateSet(ctx context.Context, input SetInput) (*Set, error) {
	set := new(Set)
	_, iset, err := r.Repository.CreateSet(ctx, input.Members)
	if err != nil {
		return set, err
	}
	set.Members = input.Members
	for _, v := range iset {
		set.IntersectingSets = append(set.IntersectingSets, Set{
			Members:          v,
			IntersectingSets: []Set{},
		})
	}
	return set, nil
}

func (r *queryResolver) Sets(ctx context.Context) ([]Set, error) {
	var setSlice []Set
	smap, imap, err := r.Repository.ListSetsWithIntersectingSets(ctx)
	if err != nil {
		return setSlice, err
	}
	for k, v := range smap {
		var s Set
		s.Members = v
		for _, isid := range imap[k] {
			s.IntersectingSets = append(s.IntersectingSets, Set{
				Members:          smap[isid],
				IntersectingSets: []Set{},
			})
		}
		setSlice = append(setSlice, s)
	}
	return setSlice, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
