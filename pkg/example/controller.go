package example

import (
	"log"
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/fields"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type Controller struct {
	jamServerClient *rest.RESTClient
	clientset       *kubernetes.Clientset

	jamServerInformer cache.SharedIndexInformer
}

func newController(jamServerClient *rest.RESTClient, clientset *kubernetes.Clientset) *Controller {
	controller := &Controller{
		jamServerClient: jamServerClient,
		clientset:       clientset,
	}

	jamServerInformer := cache.NewSharedIndexInformer(
		cache.NewListWatchFromClient(jamServerClient, "jamservers", api.NamespaceAll, fields.Everything()),
		&JamServer{},
		resyncPeriod,
		cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	)

	jamServerInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		StartFunc: controller.jamServerStart,
		StopFunc:  controller.jamServerStop,
	})

	controller.jamServerInformer = jamServerInformer

	return controller
}

func (z *Controller) Run(stopCh <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	wg.Add(1)

	if err := EnsureJamServerThirdPartyResource(z.clientset); err != nil {
		log.Fatalf("Couldn't create ThirdPartyResource: %s", err)
	}

	go z.jamServerInformer.Run(stopCh)

	<-stopCh
}

func (z *Controller) jamServerStart(obj interface{}) {
	jamServer := obj.(*JamServer)
	log.Printf("JamServer START: %s said %s... :)", jamServer.Spec.Owner, jamServer.Spec.Message)
}

func (z *Controller) jamServerStop(obj interface{}) {
	jamServer := obj.(*JamServer)
	log.Printf("JamServer STOP: %s said %s... :(", jamServer.Spec.Owner, jamServer.Spec.Message)
}
