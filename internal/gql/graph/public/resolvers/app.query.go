package pubresolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"github.com/azarc-io/example-app/internal/gql/graph/common/model"
	pubgraph "github.com/azarc-io/example-app/internal/gql/graph/public"
)

// GetValue is the resolver for the getValue field.
func (r *queryResolver) GetValue(ctx context.Context) (*model.ValueOutput, error) {
	return r.AppService.GQLGetValue(ctx)
}

// Query returns pubgraph.QueryResolver implementation.
func (r *Resolver) Query() pubgraph.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
