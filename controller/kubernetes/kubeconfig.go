package kubernetes

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/torchiaf/Sensors/controller/config"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func initKubeconfig() (*kubernetes.Clientset, client.Client, *rest.Config) {
	var restConfig *rest.Config
	var err error

	if config.Config.IsDev {
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			fmt.Print(filepath.Join(home, "work/Sensors/kubeconfig.yaml"))
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, "work/Sensors/kubeconfig.yaml"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		restConfig, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err)
		}
	} else {
		// creates the in-cluster config
		restConfig, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		panic(err.Error())
	}

	runtimeClient, err := client.New(restConfig, client.Options{})
	if err != nil {
		panic(err.Error())
	}

	return clientset, runtimeClient, restConfig
}

var Clientset, RuntimeClient, RestConfig = initKubeconfig()
