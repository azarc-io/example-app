package event

import (
	"context"
	"github.com/azarc-io/example-app/internal/types"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

/************************************************************************/
/* TYPES
/************************************************************************/

type (
	appEventManager struct {
		opts *types.Options
		log  zerolog.Logger
	}
)

/************************************************************************/
/* VALUE EVENT
/************************************************************************/

// SubscribeToValueChanges subscribes to value changed events and streams the updates through an event channel
func (a *appEventManager) SubscribeToValueChanges(_ context.Context) (eventCh chan string, close func(), err error) {
	var (
		nc  = a.opts.Service.Nats().Client()
		sub *nats.Subscription
	)

	eventCh = make(chan string)

	sub, err = nc.Subscribe(types.ValueChangedTopic, func(msg *nats.Msg) {
		eventCh <- string(msg.Data)
	})

	return eventCh, func() {
		if sub != nil {
			if err := sub.Unsubscribe(); err != nil {
				a.log.Warn().Err(err).Msgf("error while ubsubcribing")
			}
		}
	}, err
}

// DispatchValueChangedEvent dispatches a value changed event
func (a *appEventManager) DispatchValueChangedEvent(ctx context.Context, value string) error {
	nc := a.opts.Service.Nats().Client()
	return nc.Publish(types.ValueChangedTopic, []byte(value))
}

/************************************************************************/
/* FACTORY
/************************************************************************/

func NewAppEventManager(opts *types.Options, log zerolog.Logger) types.AppEventManager {
	return &appEventManager{
		opts: opts,
		log:  log,
	}
}
