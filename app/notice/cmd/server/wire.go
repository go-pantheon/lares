//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/google/wire"
	"github.com/luffy050596/rec-account/app/notice/internal/admin"
	"github.com/luffy050596/rec-account/app/notice/internal/conf"
	"github.com/luffy050596/rec-account/app/notice/internal/data"
	"github.com/luffy050596/rec-account/app/notice/internal/http"
	"github.com/luffy050596/rec-account/app/notice/internal/server"
	"github.com/luffy050596/rec-net/health"
)

func initApp(*conf.Server, *conf.Label, *conf.Registry, *conf.Data, log.Logger, *health.Server) (*kratos.App, func(), error) {
	panic(wire.Build(data.ProviderSet, server.ProviderSet, http.ProviderSet, admin.ProviderSet, newApp))
}
