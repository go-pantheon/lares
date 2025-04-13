//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-net/health"
	"github.com/go-pantheon/lares/app/notice/internal/admin"
	"github.com/go-pantheon/lares/app/notice/internal/conf"
	"github.com/go-pantheon/lares/app/notice/internal/data"
	"github.com/go-pantheon/lares/app/notice/internal/http"
	"github.com/go-pantheon/lares/app/notice/internal/server"
	"github.com/google/wire"
)

func initApp(*conf.Server, *conf.Label, *conf.Registry, *conf.Data, log.Logger, *health.Server) (*kratos.App, func(), error) {
	panic(wire.Build(data.ProviderSet, server.ProviderSet, http.ProviderSet, admin.ProviderSet, newApp))
}
