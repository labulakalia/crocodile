package registry

import (
	"fmt"
	"github.com/labulaka521/logging"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-plugins/registry/etcdv3"
)

func Etcd(addrs ...string) registry.Registry {
	reg := etcdv3.NewRegistry(
		registry.Addrs(addrs...),
	)
	return reg
}

func GetEtcdListServicesIP(srvName string, etcdRegAddrs ...string) (regips []string, err error) {
	var (
		services []*registry.Service
		service  *registry.Service
		node     *registry.Node
	)

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

func Consul(addrs ...string) registry.Registry {
	reg := consul.NewRegistry(
		registry.Addrs(addrs...),
	)
	return reg
}
