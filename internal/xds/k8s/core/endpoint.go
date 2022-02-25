package core

import (
	"github.com/kanengo/egoist/internal/xds"
	"sync"
)

func (em *EndpointManager) Watch(resourceNames []string) []<-chan []*xds.DiscoveryResponse {
	eps := make([]*Endpoint, 0, len(resourceNames))
	em.mu.Lock()
	for _, resourceName := range resourceNames {
		if res, ok := em.endpoints[resourceName]; ok {
			eps = append(eps, res)
		} else {
			res = &Endpoint{
				ResourceName: resourceName,
				watchers:     nil,
				mu:           sync.Mutex{},
			}
			em.endpoints[resourceName] = res
			eps = append(eps, res)
		}
	}
	em.mu.Unlock()
	chs := make([]<-chan []*xds.DiscoveryResponse, 0, len(eps))
	for _, ep := range eps {
		chs = append(chs, ep.Watch())
	}

	return chs
}

func (e *Endpoint) Watch() <-chan []*xds.DiscoveryResponse {
	e.mu.Lock()
	defer e.mu.Unlock()
	ch := make(chan []*xds.DiscoveryResponse, 8)
	e.watchers = append(e.watchers, ch)
	return ch
}
