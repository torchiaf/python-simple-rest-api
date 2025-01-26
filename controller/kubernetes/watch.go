package kubernetes

import (
	"fmt"
	"log"

	cm "github.com/torchiaf/Sensors/controller/circuit"
	"github.com/torchiaf/Sensors/controller/config"
	"github.com/torchiaf/Sensors/controller/models"
	"github.com/torchiaf/Sensors/controller/utils"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
)

var release = config.Config.Release

func watch(resource string, addFunc func(interface{})) {

	clusterClient, err := dynamic.NewForConfig(RestConfig)
	if err != nil {
		log.Fatalln(err)
	}

	fac := dynamicinformer.NewFilteredDynamicSharedInformerFactory(
		clusterClient,
		0,
		release.Namespace,
		nil,
	)

	informer := fac.ForResource(schema.GroupVersionResource{
		Group:    release.Group,
		Version:  release.Version,
		Resource: resource,
	}).Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: addFunc,
		UpdateFunc: func(oldObj, newObj interface{}) {
			// TODO
			fmt.Print("upd circuit", newObj)
		},
		DeleteFunc: func(obj interface{}) {
			// TODO
			fmt.Print("delete circuit", obj)
		},
	})

	informer.Run(make(chan struct{}))
}

func circuitsHandler(obj interface{}) {
	circuit := utils.ObjToStruct[*models.Circuit](obj)

	circuitManager := cm.New(Clientset, release)

	circuitManager.Create(circuit)
}

func WatchResources() {
	watch("circuits", circuitsHandler)
}
