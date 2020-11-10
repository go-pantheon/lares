//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/google/wire"
	"github.com/luffy050596/rec-account/app/account/internal/admin"
	"github.com/luffy050596/rec-account/app/account/internal/conf"
	"github.com/luffy050596/rec-account/app/account/internal/data"
	"github.com/luffy050596/rec-account/app/account/internal/http"
	"github.com/luffy050596/rec-account/app/account/internal/server"
	"github.com/luffy050596/rec-net/health"
)

func initApp(*conf.Server, *conf.Label, *conf.Registry, *conf.Data, *conf.Platform, log.Logger, *health.Server) (*kratos.App, func(), error) {
	panic(wire.Build(data.ProviderSet, server.ProviderSet, http.ProviderSet, newApp))
}
