package cluster

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ClusterLocation = string

const (
	IN_CLUSTER  ClusterLocation = "IN"
	OFF_CLUSTER ClusterLocation = "OUT"
)

type Cluster struct {
	location   ClusterLocation
	kubeconfig string
	Client     *kubernetes.Clientset
}

func (c *Cluster) inClusterClient() *kubernetes.Clientset {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the client
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return client
}

func (c *Cluster) offClusterClient(kubeconfig string) *kubernetes.Clientset {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the client
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return client
}

func New(location ClusterLocation, kubeconfig string) *Cluster {
	var client *kubernetes.Clientset

	c := &Cluster{
		location:   location,
		kubeconfig: kubeconfig,
		Client:     client,
	}
	switch location {
	case IN_CLUSTER:
		client = c.inClusterClient()
	case OFF_CLUSTER:
		client = c.offClusterClient(kubeconfig)
	}

	return c
}

func DefaultCLuster() *Cluster {
	return New(IN_CLUSTER, "")
}
