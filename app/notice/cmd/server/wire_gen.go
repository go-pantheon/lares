// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/luffy050596/rec-net/health"
	biz2 "github.com/luffy050596/rec-account/app/notice/internal/admin/biz"
	data3 "github.com/luffy050596/rec-account/app/notice/internal/admin/data"
	v1_2 "github.com/luffy050596/rec-account/app/notice/internal/admin/service/v1"
	"github.com/luffy050596/rec-account/app/notice/internal/conf"
	"github.com/luffy050596/rec-account/app/notice/internal/data"
	"github.com/luffy050596/rec-account/app/notice/internal/http/biz"
	data2 "github.com/luffy050596/rec-account/app/notice/internal/http/data"
	"github.com/luffy050596/rec-account/app/notice/internal/http/service/v1"
	"github.com/luffy050596/rec-account/app/notice/internal/server"
)

// Injectors from wire.go:

func initApp(confServer *conf.Server, label *conf.Label, registry *conf.Registry, confData *conf.Data, logger log.Logger, healthServer *health.Server) (*kratos.App, func(), error) {
	dataData, cleanup, err := data.NewData(confData)
	if err != nil {
		return nil, nil, err
	}
	noticeRepo, err := data2.NewNoticeData(dataData, logger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	noticeUseCase := biz.NewNoticeUseCase(noticeRepo, logger)
	noticeInterface := v1.NewNoticeInterface(logger, noticeUseCase)
	bizNoticeRepo, err := data3.NewNoticeData(dataData, logger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	bizNoticeUseCase := biz2.NewNoticeUseCase(bizNoticeRepo, logger)
	noticeAdmin := v1_2.NewNoticeAdmin(logger, bizNoticeUseCase)
	httpServer := server.NewHTTPServer(label, confServer, logger, noticeInterface, noticeAdmin)
	grpcServer := server.NewGRPCServer(confServer, logger, noticeAdmin)
	registrar, err := server.NewRegistrar(registry)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	app := newApp(logger, httpServer, grpcServer, healthServer, label, registrar)
	return app, func() {
		cleanup()
	}, nil
}
