package types

import (
	"context"
	"github.com/azarc-io/example-app/internal/gql/graph/common/model"
)

type (
	// App contract for the app entry point
	App interface {
		// PreStart called before the service is started, register http endpoints etc. here
		PreStart() error
		// PostStart called after the service has started, run any background processes etc. here
		PostStart() error
		// PreStop called before the service shuts down, use this to gracefully stop your service
		PreStop() error
	}

	// AppService contract for graphql handlers and other app specific business logic
	AppService interface {
		GQLSetValue(ctx context.Context, input model.ValueInput) (*model.ValueOutput, error)
		GQLGetValue(ctx context.Context) (*model.ValueOutput, error)
		GQLSubscribeToValueChanges(ctx context.Context) (<-chan *model.ValueOutput, error)
	}

	// AppEventManager handles applications events
	AppEventManager interface {
		DispatchValueChangedEvent(ctx context.Context, value string) error
		SubscribeToValueChanges(ctx context.Context) (eventCh chan string, close func(), err error)
	}
)
