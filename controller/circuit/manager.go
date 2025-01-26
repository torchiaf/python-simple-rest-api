package circuit

import (
	"bytes"
	"context"
	"fmt"

	"github.com/torchiaf/Sensors/controller/models"
	"github.com/torchiaf/Sensors/controller/utils"

	corev1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func compressCode(path string) ([]byte, error) {
	var buf bytes.Buffer
	err := utils.Compress(path, &buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func fillConfigMap(name string, releaseName string, circuit *models.Circuit, source []byte) *v1.ConfigMap {
	configMap := utils.ParseK8sResource[*v1.ConfigMap]("sensors/circuit/configmap.yaml")

	configMap.Name = name
	configMap.Labels["app.kubernetes.io/name"] = releaseName
	configMap.Labels["app.kubernetes.io/instance"] = releaseName
	configMap.Annotations["meta.helm.sh/release-name"] = releaseName
	configMap.Annotations["meta.helm.sh/release-namespace"] = releaseName

	binaryData := make(map[string][]byte)

	// TODO source-code should be a const
	binaryData["source-code.gz"] = source

	configMap.BinaryData = binaryData

	return configMap
}

func fillDeployment(name string, releaseName string, circuit *models.Circuit) *corev1.Deployment {
	deployment := utils.ParseK8sResource[*corev1.Deployment]("sensors/circuit/deployment.yaml")

	deployment.Name = name
	deployment.Labels["app.kubernetes.io/name"] = releaseName
	deployment.Labels["app.kubernetes.io/instance"] = releaseName
	deployment.Annotations["meta.helm.sh/release-name"] = releaseName
	deployment.Annotations["meta.helm.sh/release-namespace"] = releaseName

	deployment.Spec.Selector.MatchLabels["app.kubernetes.io/name"] = releaseName
	deployment.Spec.Selector.MatchLabels["app.kubernetes.io/instance"] = releaseName

	deployment.Spec.Template.Labels["app.kubernetes.io/name"] = releaseName
	deployment.Spec.Template.Labels["app.kubernetes.io/instance"] = releaseName

	deployment.Spec.Template.Spec.ServiceAccountName = releaseName

	deployment.Spec.Template.Spec.Containers[0].EnvFrom[0].SecretRef.Name = fmt.Sprintf("%s-%s", releaseName, "rabbitmq")

	deployment.Spec.Template.Spec.Volumes[0].ConfigMap.Name = releaseName                                              // modules.yaml
	deployment.Spec.Template.Spec.Volumes[1].ConfigMap.Name = fmt.Sprintf("%s-%s", releaseName, circuit.Metadata.Name) // circuit's source code

	return deployment
}

func (cm *CircuitManager) Create(circuit *models.Circuit) error {

	circuitName := fmt.Sprintf("%s-%s", cm.release.Name, circuit.Metadata.Name)

	// TODO dest should be a param, remote test for prod
	source, err := compressCode("sensors/test-compress")
	if err != nil {
		return fmt.Errorf("error compressing circuit source code %q", err)
	}

	_ = cm.clientset.
		CoreV1().
		ConfigMaps(cm.release.Namespace).
		Delete(
			context.Background(),
			circuitName,
			metav1.DeleteOptions{},
		)

	// Create the circuit's configMap from template
	configMap := fillConfigMap(circuitName, cm.release.Name, circuit, source)
	_, err = cm.clientset.
		CoreV1().
		ConfigMaps(cm.release.Namespace).
		Create(
			context.Background(),
			configMap,
			metav1.CreateOptions{},
		)
	if err != nil {
		return fmt.Errorf("error creating configmap %q", err)
	}

	var gp = int64(0)
	_ = cm.clientset.
		AppsV1().
		Deployments(cm.release.Namespace).
		Delete(
			context.Background(),
			circuitName,
			metav1.DeleteOptions{
				GracePeriodSeconds: &gp,
			},
		)

	// Create the circuit's deployment from template
	deployment := fillDeployment(circuitName, cm.release.Name, circuit)
	_, err = cm.clientset.
		AppsV1().
		Deployments(cm.release.Namespace).
		Create(
			context.Background(),
			deployment,
			metav1.CreateOptions{},
		)
	if err != nil {
		return fmt.Errorf("error creating deployment %q", err)
	}

	fmt.Printf("Created circuit %s", circuit.Metadata.Name)

	return nil
}

type CircuitManager struct {
	clientset *kubernetes.Clientset
	release   models.Release
}

func New(clientset *kubernetes.Clientset, release models.Release) *CircuitManager {
	return &CircuitManager{
		clientset: clientset,
		release:   release,
	}
}
