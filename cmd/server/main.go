package main

import (
	"flag"
	"os"

	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2"
	kratoslog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gomicroim/gomicroim/v2/pkg/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs/config.yaml", "config path, eg: -conf config.yaml")
}

func newApp(logger *log.Logger, gs *grpc.Server, hs *http.Server, registry *etcd.Registry) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
		),
		kratos.Registrar(registry),
	)
}

func main() {
	flag.Parse()
	kratoslog.SetLogger(log.MustNewLogger(id, Name, Version, true, 2))
	log.SetGlobalLogger(log.MustNewLogger(id, Name, Version, true, 0))

	bc := conf.MustLoad(flagconf)

	// register etcd
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: bc.Registry.Etcd.Endpoints,
	})
	if err != nil {
		panic(err)
	}
	log.L.Info("register etcd",
		zap.Strings("endpoints", bc.Registry.Etcd.Endpoints))
	reg := etcd.New(etcdClient)

	app, cleanup, err := wireApp(bc.Server, bc.Data,
		log.MustNewLogger(id, Name, Version, true, 4), // fix kratos caller stack
		log.L, reg)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
