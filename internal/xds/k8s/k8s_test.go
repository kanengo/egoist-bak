package k8s

import (
	"fmt"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/labels"

	"k8s.io/client-go/informers"
)

func TestNewClient(t *testing.T) {
	c := NewClient("")
	//fmt.Println(c.CoreV1().Pods("test").List(context.TODO(), v1.ListOptions{}))
	informerFactory := informers.NewSharedInformerFactoryWithOptions(c.Clientset, time.Hour)
	stopper := make(chan struct{})
	informerFactory.Core().V1().Pods().Informer()
	go informerFactory.Start(stopper)
	informerFactory.WaitForCacheSync(stopper)
	time.Sleep(3 * time.Second)
	list, err := informerFactory.Core().V1().Pods().Lister().Pods("test").List(labels.Everything())
	fmt.Println(list, err)
}
