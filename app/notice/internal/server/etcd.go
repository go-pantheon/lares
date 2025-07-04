package server

import (
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-pantheon/lares/app/notice/internal/conf"
	"github.com/google/wire"
	"github.com/pkg/errors"
	etcdclient "go.etcd.io/etcd/client/v3"
)

var ProviderSet = wire.NewSet(NewHTTPServer, NewGRPCServer, NewRegistrar)

func NewRegistrar(conf *conf.Registry) (registry.Registrar, error) {
	client, err := etcdclient.New(etcdclient.Config{
		Endpoints: conf.Etcd.Endpoints,
		Username:  conf.Etcd.Username,
		Password:  conf.Etcd.Password,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "[etcdclient.New] etcd client create failed")
	}

	return etcd.New(client), nil
}
