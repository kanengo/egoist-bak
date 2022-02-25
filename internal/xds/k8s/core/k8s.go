package core

import (
	"errors"
	"fmt"
	"github.com/kanengo/egoist/internal/xds"
	"go.uber.org/atomic"
	"path/filepath"
	"time"

	"k8s.io/client-go/informers"

	v1 "k8s.io/api/core/v1"

	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/homedir"
)

type K8sDiscovery struct {
	informerFactory informers.SharedInformerFactory
	stopCh          chan struct{}
	endpointManger  *EndpointManager

	closed atomic.Bool
}

var (
	ErrClientHadClosed = errors.New("k8s client had been closed")
)

func NewClient(kubeconfig string, namespace string, defaultSyncPeriod time.Duration) *K8sDiscovery {
	var err error
	config, err := rest.InClusterConfig()
	if err != nil {
		if err != rest.ErrNotInCluster {
			panic(err.Error())
		}
		//not in cluster
		if kubeconfig == "" {
			if home := homedir.HomeDir(); home != "" {
				kubeconfig = filepath.Join(home, ".kube", "config")
			}
		}
	} else {
		kubeconfig = ""
	}

	var (
		clientSet *kubernetes.Clientset
	)

	if config == nil {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}

	clientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	c := &K8sDiscovery{
		stopCh: make(chan struct{}),
	}

	informerFactory := informers.NewSharedInformerFactoryWithOptions(clientSet, defaultSyncPeriod, informers.WithNamespace(namespace))
	eh := &eventHandler{}
	informerFactory.Core().V1().Pods().Informer().AddEventHandler(eh)
	informerFactory.Start(c.stopCh)
	informerFactory.WaitForCacheSync(c.stopCh)

	c.informerFactory = informerFactory

	return c
}

func (d *K8sDiscovery) Watch(req *xds.DiscoveryRequest) ([]<-chan []*xds.DiscoveryResponse, error) {
	if d.closed.Load() {
		return nil, ErrClientHadClosed
	}
	switch req.TypeUrl {
	case xds.EDS:
		return d.endpointManger.Watch(req.ResourceNames), nil
	default:

	}

	return nil, fmt.Errorf("watch a invalid resource type:%s", req.TypeUrl)
}

func (d *K8sDiscovery) Close() {
	if !d.closed.CAS(false, true) {
		return
	}
	close(d.stopCh)

}

type eventHandler struct {
}

func (in *eventHandler) OnAdd(obj interface{}) {
	switch res := obj.(type) {
	case *v1.Pod:
		pod := Pod{res}
		fmt.Println("ONAdd", pod.Name, pod.Ready())
	default:
	}

}

func (in *eventHandler) OnUpdate(oldObj, newObj interface{}) {
	switch oldRes := oldObj.(type) {
	case *v1.Pod:
		newRes := newObj.(*v1.Pod)
		oldOne, newOne := Pod{oldRes}, Pod{newRes}
		if oldOne.Ready() != newOne.Ready() {

		}
	default:
	}

}

func (in *eventHandler) OnDelete(obj interface{}) {
	switch res := obj.(type) {
	case *v1.Pod:
		pod := Pod{res}
		fmt.Println("onDelete", pod.Name)
	default:
	}
}
