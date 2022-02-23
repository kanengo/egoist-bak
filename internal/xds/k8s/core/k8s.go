package core

import (
	"fmt"
	"path/filepath"

	v1 "k8s.io/api/core/v1"

	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/homedir"
)

type Client struct {
	*kubernetes.Clientset
}

func NewClient(kubeconfig string) *Client {
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
		clientset *kubernetes.Clientset
	)

	if config == nil {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	c := &Client{
		Clientset: clientset,
	}

	return c
}

type eventHandler struct {
}

func (in *eventHandler) OnAdd(obj interface{}) {
	pod := obj.(*v1.Pod)
	pp := Pod{pod}
	fmt.Println("ONAdd", pp.Name, pp.Ready())
}

func (in *eventHandler) OnUpdate(oldObj, newObj interface{}) {
	oldOne, newOne := Pod{oldObj.(*v1.Pod)}, Pod{newObj.(*v1.Pod)}
	fmt.Println("OnUpdate", oldOne.Name, oldOne.Ready(), newOne.Name, newOne.Ready())
}

func (in *eventHandler) OnDelete(obj interface{}) {
	pod := obj.(*v1.Pod)
	pp := Pod{pod}
	fmt.Println("OnDelete", pp.Name, pp.Ready())
}
