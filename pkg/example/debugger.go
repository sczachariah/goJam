package example

import (
	"log"
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/runtime"
	"k8s.io/client-go/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type Debugger struct {
	podInformer cache.SharedIndexInformer
}

func newDebugger(client *kubernetes.Clientset) *Debugger {
	debugger := &Debugger{}

	podInformer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options api.ListOptions) (runtime.Object, error) {
				return client.Core().Pods(v1.NamespaceAll).List(options)
			},
			WatchFunc: func(options api.ListOptions) (watch.Interface, error) {
				return client.Core().Pods(v1.NamespaceAll).Watch(options)
			},
		},
		&v1.Pod{},
		resyncPeriod,
		cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	)

	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    debugger.podAdd,
		UpdateFunc: debugger.podUpdate,
		DeleteFunc: debugger.podDelete,
	})

	debugger.podInformer = podInformer

	return debugger
}

func (d *Debugger) Run(stopCh <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	wg.Add(1)

	go d.podInformer.Run(stopCh)

	<-stopCh
}

func (d *Debugger) podAdd(obj interface{}) {
	pod := obj.(*v1.Pod)
	log.Printf("ADD %s/%s", pod.Namespace, pod.Name)
}

func (d *Debugger) podDelete(obj interface{}) {
	pod := obj.(*v1.Pod)
	log.Printf("DELETE %s/%s", pod.Namespace, pod.Name)
}

func (d *Debugger) podUpdate(cur, old interface{}) {
	pod := cur.(*v1.Pod)
	log.Printf("UPDATE %s/%s", pod.Namespace, pod.Name)
}
