package apimachinery

import (
	errors "github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func hasEphemeralContainer(ecs []v1.EphemeralContainer, ecName string) bool {
	for _, ec := range ecs {
		if ec.Name == ecName {
			return true
		}
	}
	return false
}

func createDebugContainerObject(name string) v1.EphemeralContainer {
	return v1.EphemeralContainer{
		EphemeralContainerCommon: v1.EphemeralContainerCommon{
			Name:                     name,
			Image:                    "nicolaka/netshoot", // TODO: find a better place to define the image
			ImagePullPolicy:          v1.PullIfNotPresent,
			TerminationMessagePolicy: v1.TerminationMessageReadFile,
			Stdin:                    true,
			//SecurityContext: &v1.SecurityContext{
			//	Capabilities: &v1.Capabilities{
			//		Add: []v1.Capability{"NET_ADMIN"},
			//	},
			//},
		},
	}
}

func CreateOrUpdateEphemeralContainer(config *rest.Config, namespace, podName, ephemeralContainerName string) error {
	debugContainer := createDebugContainerObject(ephemeralContainerName)
	client, err := clientset.NewForConfig(config)
	if err != nil {
		return err
	}

	pods := client.CoreV1().Pods(namespace)
	ec, err := pods.GetEphemeralContainers(podName, metav1.GetOptions{})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return errors.Wrap(err, "ephemeral containers are not enabled for this cluster")
		}
		return err
	}

	if !hasEphemeralContainer(ec.EphemeralContainers, debugContainer.Name) {
		ec.EphemeralContainers = append(ec.EphemeralContainers, debugContainer)
		_, err = pods.UpdateEphemeralContainers(podName, ec)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreatePod(config *rest.Config, name string, image string, args []string, serviceAccountName string) error {

	// Discuss : Can we use a single client here
	client := clientset.NewForConfigOrDie(config)
	pod := &v1.Pod{}
	pod.Name = name
	container := v1.Container{Name: name,
		Image:           image,
		ImagePullPolicy: v1.PullAlways,
		Args:            args,
	}
	pod.Spec.Containers = append(pod.Spec.Containers, container)
	pod.Spec.RestartPolicy = "OnFailure"
	pod.Spec.ServiceAccountName = serviceAccountName

	_, err := client.CoreV1().Pods(v1.NamespaceDefault).Create(pod)
	return err

}

func GetPod(config *rest.Config, name string) (*v1.Pod, error) {

	client := clientset.NewForConfigOrDie(config)
	pod, err := client.CoreV1().Pods(v1.NamespaceDefault).Get(name, metav1.GetOptions{})
	return pod, err

}

func DeletePod(config *rest.Config, name string) error {

	client := clientset.NewForConfigOrDie(config)
	err := client.CoreV1().Pods(v1.NamespaceDefault).Delete(name, &metav1.DeleteOptions{})
	return err

}

func GetLogs(config *rest.Config, name string) *rest.Request {

	client := clientset.NewForConfigOrDie(config)
	data := client.CoreV1().Pods(v1.NamespaceDefault).GetLogs(name, &v1.PodLogOptions{})
	return data

}
