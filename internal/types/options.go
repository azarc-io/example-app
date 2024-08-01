package types

import "github.com/azarc-io/verathread-next-common/service"

type (
	Options struct {
		Config  *Config
		Service *service.AppService
	}

	Option func(o *Options)
)

func WithConfig(cfg *Config) Option {
	return func(o *Options) {
		o.Config = cfg
	}
}

func WithService(svc *service.AppService) Option {
	return func(o *Options) {
		o.Service = svc
	}
}
