package xds

type Client interface {
	GetResources(response *DiscoveryRequest) (*DiscoveryResponse, error)
	WatchResource(response *DiscoveryRequest) ([]chan *DiscoveryResponse, error)
	Clone() Client
	Close()
}
