package pvtresolvers

import "github.com/azarc-io/example-app/internal/types"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Opts       *types.Options
	AppService types.AppService
}
