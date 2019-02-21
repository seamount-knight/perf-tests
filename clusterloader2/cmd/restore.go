package main

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/perf-tests/clusterloader2/pkg/flags"
	"k8s.io/perf-tests/clusterloader2/pkg/framework"
)

func checkNode() {
	kubeconfigPath := ""
	flags.StringEnvVar(&kubeconfigPath, "kubeconfig", "KUBECONFIG", "", "Path to the kubeconfig file")

	// conf, err := config.PrepareConfig(kubeconfigPath)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// k8sClient, err := clientset.NewForConfig(conf)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	clientSet, err := framework.NewMultiClientSet(kubeconfigPath, 1)
	if err != nil {
		fmt.Println(err)
		return
	}

	node, err := clientSet.GetClient().CoreV1().Nodes().Get("", metav1.GetOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	node.Status
}
