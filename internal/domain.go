package app

import (
	"github.com/azarc-io/example-app/internal/event"
	pvtgraph "github.com/azarc-io/example-app/internal/gql/graph/private"
	pvtresolvers "github.com/azarc-io/example-app/internal/gql/graph/private/resolvers"
	pubgraph "github.com/azarc-io/example-app/internal/gql/graph/public"
	pubresolvers "github.com/azarc-io/example-app/internal/gql/graph/public/resolvers"
	"github.com/azarc-io/example-app/internal/service"
	"github.com/azarc-io/example-app/internal/types"
	app2 "github.com/azarc-io/verathread-next-common/common/app"
	appuc "github.com/azarc-io/verathread-next-common/usecase/app"
	graphqluc "github.com/azarc-io/verathread-next-common/usecase/graphql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"strings"
	"time"
)

/************************************************************************/
/* TYPES
/************************************************************************/

type (
	app struct {
		opts            *types.Options
		publicAPI       graphqluc.GraphQLUseCase
		privateAPI      graphqluc.GraphQLUseCase
		log             zerolog.Logger
		appService      types.AppService
		registration    appuc.AppUseCase
		appEventManager types.AppEventManager
	}
)

/************************************************************************/
/* LIFECYCLE
/************************************************************************/

// PreStart called before the service is started, register http endpoints etc. here
func (a *app) PreStart() error {
	// initialize and register services
	if err := a.registerServicesAndEventManagers(); err != nil {
		return err
	}

	// register graphql end points
	if err := a.registerGqlAPI(); err != nil {
		return err
	}

	// register web app router
	a.registerWebAppRoute()

	return nil
}

// PostStart called after the service has started, run any background processes etc. here
func (a *app) PostStart() error {
	// run leadership routines
	a.runOnLeader()

	return a.registerApplication()
}

// PreStop called before the service shuts down, use this to gracefully stop your service
func (a *app) PreStop() error {
	return nil
}

/************************************************************************/
/* LEADERSHIP
/************************************************************************/

// runOnLeader watches for changes in leadership and spawns or de-spawns routines based on being leader or follower
func (a *app) runOnLeader() {
	a.opts.Service.Redis().SubscribeToElectionEvents(func(onPromote <-chan time.Time, onDemote <-chan time.Time) {
		select {
		case <-a.opts.Service.Context().Done():
			a.log.Info().Msgf("stopping leadership routines, service is shutting down")
		case <-onPromote:
			a.log.Info().Msgf("running leadership routines")
		case <-onDemote:
			a.log.Info().Msgf("stopping leadership routines")
		}
	})
}

/************************************************************************/
/* SERVICES & EVENT DRIVEN MANAGERS
/************************************************************************/

// registerServicesAndEventManagers initialises and registers services
func (a *app) registerServicesAndEventManagers() error {
	// create app event manager to handle app specific events
	a.appEventManager = event.NewAppEventManager(a.opts, a.log)
	// create service to handle inbound requests
	a.appService = service.NewAppService(a.opts, a.log, a.appEventManager)

	return nil
}

/************************************************************************/
/* REGISTRATION
/************************************************************************/

// registerApplication registers the application with the application gateway
func (a *app) registerApplication() error {
	a.registration = appuc.NewAppUseCase(
		appuc.WithLogger(a.log),
		appuc.WithGatewayUrl(a.opts.Config.Routing.GatewayURL),
		appuc.WithAppInfo(app2.RegisterAppInput{
			Name:            a.opts.Service.ServiceName(), // the gateway will use this name to proxy e.g. /module/user/*
			Id:              a.opts.Service.ServiceID(),
			Package:         "vth:azarc:" + a.opts.Service.ServiceName(),
			Version:         a.opts.Service.Version(),
			ApiUrl:          a.opts.Config.Routing.APIURL,
			WebUrl:          a.opts.Config.Routing.WebURL,
			RemoteEntryFile: "remoteEntry.js", // if proxy is true then don't need url here
			Proxy:           a.opts.Config.Routing.Proxy,
			Navigation: []app2.RegisterAppNavigationInput{
				{
					Title:    a.opts.Service.Title(),
					SubTitle: a.opts.Service.Description(),
					Module: app2.RegisterAppModule{
						ExposedModule: "./App",
						ModuleName:    strings.ReplaceAll(a.opts.Service.ServiceName(), "-", "_"),
						Path:          "/" + a.opts.Service.ServiceName(),
					},
					AuthRequired: true,
					Hidden:       false,
					Proxy:        true,
					Category:     app2.RegisterAppCategorySetting,
					Children:     []app2.RegisterChildAppNavigationInput{},
				},
			},
		}),
	)

	return a.registration.Start()
}

/************************************************************************/
/* API
/************************************************************************/

// registerGqlAPI registers graphql api handler
func (a *app) registerGqlAPI() error {
	// public api
	a.publicAPI = graphqluc.NewGraphQLUseCase(
		graphqluc.WithLogger(a.log),
		graphqluc.WithHTTPUseCase(a.opts.Service.PublicHTTP()),
		graphqluc.WithExecutableSchema(pubgraph.NewExecutableSchema(pubgraph.Config{
			Resolvers: &pubresolvers.Resolver{
				Opts:       a.opts,
				AppService: a.appService,
			},
		})),
	)

	// private api
	a.privateAPI = graphqluc.NewGraphQLUseCase(
		graphqluc.WithLogger(a.log),
		graphqluc.WithHTTPUseCase(a.opts.Service.PrivateHTTP()),
		graphqluc.WithExecutableSchema(pvtgraph.NewExecutableSchema(pvtgraph.Config{
			Resolvers: &pvtresolvers.Resolver{
				Opts:       a.opts,
				AppService: a.appService,
			},
		})),
	)

	return nil
}

/************************************************************************/
/* WEB APP
/************************************************************************/

// registerWebAppRoute serves the web application
func (a *app) registerWebAppRoute() {
	e := a.opts.Service.PublicHTTP().Server()
	g1 := e.Group("")
	g1.Use(
		middleware.GzipWithConfig(middleware.GzipConfig{
			Skipper: func(c echo.Context) bool {
				return !strings.Contains(c.Path(), ".html")
			},
		}),
		middleware.StaticWithConfig(middleware.StaticConfig{
			Root:  a.opts.Config.WebDir,
			Index: "index.html",
			HTML5: true,
			Skipper: func(e echo.Context) bool {
				return strings.HasPrefix(e.Path(), "/tmp") ||
					strings.HasPrefix(e.Path(), "/api") ||
					strings.HasPrefix(e.Path(), "/graphql") ||
					strings.HasPrefix(e.Path(), "/query") ||
					strings.HasPrefix(e.Path(), "/health")
			},
		}),
	)
}

/************************************************************************/
/* FACTORY
/************************************************************************/

func NewApp(opts ...types.Option) types.App {
	a := &app{
		opts: &types.Options{},
	}

	for _, opt := range opts {
		opt(a.opts)
	}

	a.log = a.opts.Service.Logger()

	return a
}
