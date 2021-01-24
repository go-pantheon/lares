package server

import (
	"github.com/google/wire"
	"github.com/luffy050596/rec-account/app/notice/internal/conf"
	"github.com/pkg/errors"
)

var ProviderSet = wire.NewSet(NewHTTPServer, NewGRPCServer, NewRegistrar)

func NewRegistrar(conf *conf.Registry) (registry.Registrar, error) {
	client, err := etcdclient.New(etcdclient.Config{
		Endpoints: conf.Etcd.Endpoints,
		Username:  conf.Etcd.Username,
		Password:  conf.Etcd.Password,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "[etcdclient.New] etcd 客户端创建失败。")
	}

	return etcd.New(client), nil
}
