package schedule

import (
	"context"
	"fmt"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/model"
	"github.com/labulaka521/crocodile/core/utils/define"
	"go.uber.org/zap"
	"math"
	"math/rand"
	"time"
)

// Select a next run host

// Next will return next run host
// if Next is nil,because not find valid host
type Next func() *define.Host

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GetRoutePolicy return a type Next, it will return a host
func GetRoutePolicy(hgid string, routepolicy define.RoutePolicy) Next {
	switch routepolicy {
	case define.Random:
		return random(hgid)
	case define.RoundRobin:
		return roundRobin(hgid)
	case define.Weight:
		return weight(hgid)
	case define.LeastTask:
		return leastTask(hgid)
	default:
		return defaultRoutePolicy(hgid)
	}
}

// getOnlineHosts return online worker host info
func getOnlineHosts(hgid string) ([]*define.Host, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	hg, err := model.GetHostGroupID(ctx, hgid)
	if err != nil {
		return nil, err
	}

	onlinehosts := make([]*define.Host, 0, len(hg.HostsID))
	gethosts, err := model.GetHostByIDS(ctx, hg.HostsID)
	if err != nil {
		log.Error("GetHostByIDS failed", zap.Strings("ids", hg.HostsID), zap.Error(err))
		return nil, err
	}

	for _, host := range gethosts {
		if host.Online == 0 {
			continue
		}
		if host.Stop == 1 {
			continue
		}
		onlinehosts = append(onlinehosts, host)
	}
	if len(onlinehosts) == 0 {
		err := fmt.Errorf("can not get valid host from hostgrop %s" ,hgid)
		return nil, err
	}
	return onlinehosts, nil

}


var defaultRoutePolicy = random

// random return a Next func,it will random return host
func random(hgid string) Next {
	log.Debug("add Next func Random")
	return func() *define.Host {
		hosts, err := getOnlineHosts(hgid)
		if err != nil {
			// log.Error("get host failed", zap.Error(err))
			return nil
		}
		return hosts[rand.Int()%len(hosts)]
	}
}

// roundRobin return a Next func,it will RoundRobin return host
func roundRobin(hgid string) Next {
	var i = rand.Int()
	return func() *define.Host {
		hosts, err := getOnlineHosts(hgid)
		if err != nil {
			return nil
		}
		host := hosts[i%len(hosts)]
		i++
		return host
	}
}

// weight return a Next Func,it will return host by host weight
func weight(hgid string) Next {
	return func() *define.Host {
		hosts, err := getOnlineHosts(hgid)
		if err != nil {
			return nil
		}
		allweight := 0

		for _, h := range hosts {
			allweight += h.Weight
		}
		get := rand.Int() % allweight
		pre := 0

		for _, h := range hosts {
			if pre <= get && get < pre+h.Weight {
				return h
			}
			pre += h.Weight
		}
		return nil
	}
}

// leastTask return a Next Func, it will return host by leaset host running task
func leastTask(hgid string) Next {
	return func() *define.Host {
		hosts, err := getOnlineHosts(hgid)
		if err != nil {
			return nil
		}
		// a worker max running tasks 32767
		var lasttotaltasks = int(math.MaxInt16)
		var leasetTask *define.Host
		for _, host := range hosts {
			if len(host.RunningTasks) < lasttotaltasks {
				leasetTask = host
				lasttotaltasks = len(host.RunningTasks)
			}
		}
		return leasetTask
	}
}
