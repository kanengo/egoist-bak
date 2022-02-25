package xds

// KVPair contains a key and a string.
type KVPair struct {
	Key   string
	Value string
}

type NodeMeta struct {
	Name string
	Ip   string
}

const (
	EDS = "eds"
)

type DiscoveryRequest struct {
	TypeUrl       string
	ResourceNames []string
	Node          *NodeMeta
}

type DiscoveryResponse struct {
	TypeUrl   string
	Resources []*KVPair
}
