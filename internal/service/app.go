package service

import (
	"context"
	"github.com/azarc-io/verathread-app-template/internal/gql/graph/common/model"
	"github.com/azarc-io/verathread-app-template/internal/types"
	"github.com/rs/zerolog"
)

/************************************************************************/
/* TYPES
/************************************************************************/

type appService struct {
	value string
	opts  *types.Options
	log   zerolog.Logger
	aem   types.AppEventManager
}

/************************************************************************/
/* GQL HANDLERS
/************************************************************************/

func (a *appService) GQLSetValue(ctx context.Context, input model.ValueInput) (*model.ValueOutput, error) {
	// cache the value
	a.value = input.Value
	// send out a notification that the value changed
	if err := a.aem.DispatchValueChangedEvent(ctx, input.Value); err != nil {
		return nil, err
	}

	return &model.ValueOutput{Value: a.value}, nil
}

func (a *appService) GQLGetValue(ctx context.Context) (*model.ValueOutput, error) {
	return &model.ValueOutput{Value: a.value}, nil
}

func (a *appService) GQLSubscribeToValueChanges(ctx context.Context) (<-chan *model.ValueOutput, error) {
	ch := make(chan *model.ValueOutput, 1)
	eventCh, cancel, _ := a.aem.SubscribeToValueChanges(ctx)

	go func() {
		for {
			select {
			case <-a.opts.Service.Context().Done():
				a.log.Info().Msgf("closing subscription, service is shutting down")
				close(ch)
				cancel()
				return
			case <-ctx.Done():
				a.log.Info().Msgf("closing subscription, client disconnected")
				close(ch)
				cancel()
				return
			case <-eventCh:
				ch <- &model.ValueOutput{Value: a.value}
			}
		}
	}()

	return ch, nil
}

/************************************************************************/
/* FACTORY
/************************************************************************/

func NewAppService(opts *types.Options, log zerolog.Logger, aem types.AppEventManager) types.AppService {
	return &appService{
		opts: opts,
		log:  log,
		aem:  aem,
	}
}
