package kubernetes

import (
	"context"

	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var crdLabelSelector = "sensors.io/crd-group=sensors"

func initCrdSchema() *runtime.Scheme {
	crdClientSet, err := clientset.NewForConfig(RestConfig)
	if err != nil {
		panic(err)
	}

	crdList, err := crdClientSet.ApiextensionsV1().CustomResourceDefinitions().List(
		context.Background(),
		metav1.ListOptions{
			LabelSelector: crdLabelSelector,
		})
	if err != nil {
		panic(err)
	}

	crdScheme := runtime.NewScheme()

	for _, crd := range crdList.Items {
		for _, v := range crd.Spec.Versions {
			crdScheme.AddKnownTypeWithName(
				schema.GroupVersionKind{
					Group:   crd.Spec.Group,
					Version: v.Name,
					Kind:    crd.Spec.Names.Kind,
				},
				&unstructured.Unstructured{},
			)
		}
	}

	return crdScheme
}

var CrdScheme = initCrdSchema()
