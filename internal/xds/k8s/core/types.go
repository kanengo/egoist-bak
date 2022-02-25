package core

import (
	"github.com/kanengo/egoist/internal/xds"
	v1 "k8s.io/api/core/v1"
	"sync"
)

type Pod struct {
	*v1.Pod
}

type Endpoint struct {
	ResourceName string
	watchers     []chan []*xds.DiscoveryResponse
	mu           sync.Mutex
}

type EndpointManager struct {
	endpoints map[string]*Endpoint
	mu        sync.RWMutex
}
