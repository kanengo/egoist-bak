package core

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/labels"

	"k8s.io/client-go/informers"
)

type podEventHandler struct {
}

func (p *podEventHandler) OnAdd(obj interface{}) {
	pod := obj.(*v1.Pod)
	pp := Pod{pod}
	fmt.Println("ONAdd", pp.Name, pp.Ready())
}

func (p *podEventHandler) OnUpdate(oldObj, newObj interface{}) {
	old, new1 := oldObj.(*v1.Pod), newObj.(*v1.Pod)
	pp := Pod{old}
	pp2 := Pod{new1}
	fmt.Println("OnUpdate", pp.Name, pp.Ready(), pp2.Name, pp2.Ready())
}

func (p *podEventHandler) OnDelete(obj interface{}) {
	pod := obj.(*v1.Pod)
	pp := Pod{pod}
	fmt.Println("OnDelete", pp.Name, pp.Ready())
}

func TestNewClient(t *testing.T) {
	c := NewClient("")
	//fmt.Println(c.CoreV1().Pods("test").List(context.TODO(), v1.ListOptions{}))
	informerFactory := informers.NewSharedInformerFactoryWithOptions(c.Clientset, time.Minute*5, informers.WithNamespace("test"))
	stopper := make(chan struct{})
	informerFactory.Core().V1().Pods().Informer()
	informerFactory.Start(stopper)
	informerFactory.WaitForCacheSync(stopper)
	//time.Sleep(3 * time.Second)
	handler := &podEventHandler{}
	informerFactory.Core().V1().Pods().Informer().AddEventHandler(handler)
	list, err := informerFactory.Core().V1().Pods().Lister().Pods("").List(labels.Everything())
	if err != nil {
		log.Fatalln(err)
	}
	for _, l := range list {
		p := Pod{l}
		fmt.Println("list", l.Name, p.Ready())
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)

	<-signalChan
}
