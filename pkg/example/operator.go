package example

import (
	"log"
	"sync"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	VERSION      = "0.0.0.dev"
	resyncPeriod = 5 * time.Minute
)

type Options struct {
	KubeConfig string
}

type Example struct {
	Options

	clientset       *kubernetes.Clientset
	jamServerClient *rest.RESTClient

	debugger   *Debugger
	controller *Controller
}

func New(options Options) *Example {
	config := newClientConfig(options)

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Couldn't create Kubernetes client: %s", err)
	}

	jamServerClient, err := NewJamServerClientForConfig(config)
	if err != nil {
		log.Fatalf("Couldn't create Critter client: %s", err)
	}

	example := &Example{
		Options:         options,
		clientset:       clientset,
		jamServerClient: jamServerClient,
		controller:      newController(jamServerClient, clientset),
		debugger:        newDebugger(clientset),
	}

	return example
}

func (example *Example) Run(stopCh <-chan struct{}, wg *sync.WaitGroup) {
	log.Printf("Welcome to JamServer Operator %v\n", VERSION)

	go example.controller.Run(stopCh, wg)
	go example.debugger.Run(stopCh, wg)
}

func newClientConfig(options Options) *rest.Config {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	overrides := &clientcmd.ConfigOverrides{}

	if options.KubeConfig != "" {
		rules.ExplicitPath = options.KubeConfig
	}

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides).ClientConfig()
	if err != nil {
		log.Fatalf("Couldn't get Kubernetes default config: %s", err)
	}

	return config
}
