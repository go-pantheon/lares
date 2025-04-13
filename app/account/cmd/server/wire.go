//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-net/health"
	"github.com/go-pantheon/lares/app/account/internal/admin"
	"github.com/go-pantheon/lares/app/account/internal/conf"
	"github.com/go-pantheon/lares/app/account/internal/data"
	"github.com/go-pantheon/lares/app/account/internal/http"
	"github.com/go-pantheon/lares/app/account/internal/server"
	"github.com/google/wire"
)

func initApp(*conf.Server, *conf.Label, *conf.Registry, *conf.Data, *conf.Platform, log.Logger, *health.Server) (*kratos.App, func(), error) {
	panic(wire.Build(data.ProviderSet, server.ProviderSet, http.ProviderSet, admin.ProviderSet, newApp))
}
