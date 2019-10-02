package iiii

import (
	commonError "errors"
	"flag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

const (
	defaultNamespace = "default"
)

func getHome() (string, error) {
	home := os.Getenv("HOME")
	if len(home) == 0 {
		home = os.Getenv("USERPROFILE")
		if len(home) == 0 {
			return home, commonError.New("cannot find home directory")
		}
	}
	return home, nil
}

func getKubeClientset() (*kubernetes.Clientset, error) {
	home, homeErr := getHome()
	if homeErr != nil {
		return nil, homeErr
	}
	kubeconfig := flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "")
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, commonError.New("cannot get kubeconfig")
	}
	return kubernetes.NewForConfig(config)
}
