package k8s

import (
	"time"

	"github.com/kanengo/egoist/internal/xds/k8s/core"
)

type Client struct {
	cli *core.K8sDiscovery
}

func NewXDSClient(namespace string, defaultSyncPeriod time.Duration) *Client {
	c := &Client{
		cli: core.NewClient("", namespace, defaultSyncPeriod),
	}
	return c
}
