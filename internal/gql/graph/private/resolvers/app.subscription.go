package pvtresolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"github.com/azarc-io/example-app/internal/gql/graph/common/model"
	pvtgraph "github.com/azarc-io/example-app/internal/gql/graph/private"
)

// SubscribeToValueChanges is the resolver for the subscribeToValueChanges field.
func (r *subscriptionResolver) SubscribeToValueChanges(ctx context.Context) (<-chan *model.ValueOutput, error) {
	return r.AppService.GQLSubscribeToValueChanges(ctx)
}

// Subscription returns pvtgraph.SubscriptionResolver implementation.
func (r *Resolver) Subscription() pvtgraph.SubscriptionResolver { return &subscriptionResolver{r} }

type subscriptionResolver struct{ *Resolver }
