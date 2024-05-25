package utils

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

type InvalidSetupError struct{}

func (e *InvalidSetupError) Error() string {
	return "Unable to correctly setup from configuration.  Generic message - most likely something wrong with the kubeconfig."
}

func NewKubeClientSet(useInclusterConfig bool) (*kubernetes.Clientset, error) {
	var cfg *rest.Config
	var err error

	if useInclusterConfig {

		if homedir := homedir.HomeDir(); homedir != "" {
			p := filepath.Join(homedir, ".kube", "config")

			cfg, err = clientcmd.BuildConfigFromFlags("", p)

			if err != nil {
				return nil, err
			}

		} else {
			return nil, &InvalidSetupError{}
		}
	} else {
		cfg, err = rest.InClusterConfig()

		if err != nil {
			return nil, err
		}
	}

	clientset, cerr := kubernetes.NewForConfig(cfg)

	if cerr != nil {
		return nil, err
	}

	return clientset, nil
}
