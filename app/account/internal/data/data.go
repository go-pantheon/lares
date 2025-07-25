package data

import (
	"context"

	"github.com/go-pantheon/fabrica-kit/trace/tracepg"
	"github.com/go-pantheon/fabrica-util/data/db/pg"
	"github.com/go-pantheon/fabrica-util/data/redis"
	"github.com/go-pantheon/lares/app/account/internal/conf"
	goredis "github.com/redis/go-redis/v9"
)

type Data struct {
	Pdb *pg.DB
	Rdb goredis.UniversalClient
}

func NewData(c *conf.Data) (d *Data, cleanup func(), err error) {
	var (
		pdb *pg.DB
		rdb goredis.UniversalClient
	)

	if c.Redis.Cluster {
		rdb, cleanup, err = redis.NewCluster(&goredis.ClusterOptions{
			Addrs:        []string{c.Redis.Addr},
			Password:     c.Redis.Password,
			DialTimeout:  c.Redis.DialTimeout.AsDuration(),
			WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
			ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		})
	} else {
		rdb, cleanup, err = redis.NewStandalone(&goredis.Options{
			Addr:         c.Redis.Addr,
			Password:     c.Redis.Password,
			DialTimeout:  c.Redis.DialTimeout.AsDuration(),
			WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
			ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		})
	}

	if err != nil {
		return nil, nil, err
	}

	pdb, cleanupPg, err := tracepg.NewDB(context.Background(), tracepg.DefaultPostgreSQLConfig(pg.Config{
		DSN:    c.Postgresql.Source,
		DBName: c.Postgresql.Database,
	}))
	if err != nil {
		return nil, nil, err
	}

	d = &Data{
		Pdb: pdb,
		Rdb: rdb,
	}

	combinedCleanup := func() {
		if cleanup != nil {
			cleanup()
		}

		if cleanupPg != nil {
			cleanupPg()
		}
	}

	return d, combinedCleanup, nil
}
