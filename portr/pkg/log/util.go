package log

import (
	"context"

	"github.com/origine-run/portr/pkg/cluster"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/metav1"
	"k8s.io/apimachinery/pkg/labels"
)

func New(cluster *cluster.Cluster) *Log {
	return &Log{
		cluster: cluster,
	}
}

// Run setup up log watch for every service of every namespace in the cluster
// specific namespace(s) and service(s) can be excluded
func Run(cluster *cluster.Cluster, namespaceExclusionList []string, serviceExclusionList []string) {
	l := New(cluster)
	client := l.cluster.Client

	namespaces, err := client.CoreV1().Namespaces().List(context.TODO(), metaV1.ListOptions{})
	if err != nil || len(namespaces.Items) < 1 {
		return
	}

	for _, namespace := range namespaces.Items {
		services, err := client.CoreV1().Services(namespace.Name).List(context.TODO(), metaV1.ListOptions{
			LabelSelector: labels.Set(namespace.Labels).AsSelector().String(),
		})

		if err != nil || len(services.Items) < 1 {
			continue
		}

		pods, err := client.CoreV1().Pods(namespace.Name).List(context.TODO(), metaV1.ListOptions{})
		if err != nil || len(pods.Items) < 1 {
			return
		}

		for _, pod := range pods.Items {
			containers := pod.Spec.Containers
			for _, container := range containers {
				go l.Watch(namespace.Name, pod.Name, container.Name)
			}
		}
	}
}
