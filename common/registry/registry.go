package registry

import (
	"fmt"
	"github.com/labulaka521/logging"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
	"sync"
)

var (
	Reg  registry.Registry
	once sync.Once
)

func Etcd(addrs ...string) registry.Registry {
	// 向etcd的连接不会关闭，防止创建出多个对象
	once.Do(func() {
		Reg = etcdv3.NewRegistry(
			registry.Addrs(addrs...),
		)
	})
	return Reg
}

func GetEtcdListServicesIP(srvName string, etcdRegAddrs ...string) (regips []string, err error) {
	var (
		services []*registry.Service
		service  *registry.Service
		node     *registry.Node
	)
	regips = []string{}

	if services, err = Etcd(etcdRegAddrs...).GetService(srvName); err != nil {
		logging.Errorf("GetService  %s Err: %v", srvName, err)
		err = nil
		return
	}

	for _, service = range services {
		for _, node = range service.Nodes {
			regips = append(regips, fmt.Sprintf("%s:%d", node.Address, node.Port))
		}
	}
	return
}
