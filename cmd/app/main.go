package main

import (
	app "github.com/azarc-io/verathread-app-template/internal"
	"github.com/azarc-io/verathread-app-template/internal/types"
	"github.com/azarc-io/verathread-next-common/service"
	"github.com/azarc-io/verathread-next-common/util"
	"github.com/rs/zerolog/log"
	"os"
	"runtime/debug"
)

type (
	application struct {
		svc *service.AppService
		app types.App
	}

	Config struct {
		service.Config `yaml:",inline"` // common config
		App            *types.Config    `yaml:"app"` // application custom config
	}
)

func main() {
	// graceful exit on startup panic
	defer func() {
		if err := recover(); err != nil {
			log.Error().Err(err.(error)).Msgf("service did not startup cleanly\n%s", debug.Stack())
			// cancel global context
			os.Exit(1)
		}
	}()

	// configuration
	var cfg Config
	util.PanicIfErr(service.LoadConfig(&cfg))

	a := &application{}

	a.svc = service.NewAppService(
		service.WithConfig(&cfg.Config),
		service.WithBeforeStart(func(svc *service.AppService) error {
			a.app = app.NewApp(
				types.WithConfig(cfg.App),
				types.WithService(svc),
			)
			return a.app.PreStart()
		}),
		service.WithAfterStart(func(svc *service.AppService) error {
			return a.app.PostStart()
		}),
		service.WithBeforeStop(func() error {
			return a.app.PreStop()
		}),
	)

	util.PanicIfErr(a.svc.Run())
}
