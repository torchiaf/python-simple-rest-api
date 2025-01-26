package main

import (
	"github.com/torchiaf/Sensors/controller/kubernetes"
	// "github.com/torchiaf/Sensors/controller/webserver"
)

func main() {
	// go webserver.InitWebServer()

	kubernetes.WatchResources()
}
