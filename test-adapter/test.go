package main

import (
	"fmt"

	adapter "github.com/torchiaf/Sensors/adapter"
)

// func toSelector(labels *metav1.LabelSelector) (labels.Selector, error) {
// 	return metav1.LabelSelectorAsSelector(labels)
// }

func main() {

	// matchLabels := make(map[string]string)

	// matchLabels["fleet.cattle.io/benchmark"] = "true"

	// var vvv *metav1.LabelSelector = &metav1.LabelSelector{
	// 	MatchLabels: matchLabels,
	// }

	// selector, _ := toSelector(vvv)

	// clusterLabels := make(map[string]string)
	// clusterLabels["fleet.cattle.io/benchmark"] = "true"

	// vvvvv := selector.Matches(labels.Set(clusterLabels))
	// fmt.Printf("res %v", vvvvv)

	res, _ := adapter.Read("raspberrypi-0", "dht11", "")

	fmt.Print(res)
}
