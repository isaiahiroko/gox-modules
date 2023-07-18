package event

import (
	"context"

	"github.com/origine-run/portr/pkg/cluster"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func New(cluster *cluster.Cluster) *event {
	return &event{
		cluster: *cluster,
	}
}

// list of resource in scope
var resources = []string{
	"namespaces",
	"configmaps",
	"secrets",
	"nodes",
	"pods",
	"daemonsets",
	"deployments",
	"replicasets",
	"statefulsets",
	"cronjobs",
	"jobs",
	"persistentvolumeclaims",
	"persistentvolumes",
	"storageclasses",
	"serviceaccounts",
	"services",
	"endpoints",
	"endpointslices",
	"clusterissuers",
	"issuers",
	"ingresses",
	"networkpolicies",
	"clusterrolebindings",
	"clusterroles",
	"rolebindings",
	"roles",
	"resourcequotas",
	"horizontalpodautoscalers",
	"limitranges",
	"events",
}

// Run setup up event watch for every resource of every namespace in the cluster
// specific namespace(s) and resource(s) can be excluded
func Run(cluster *cluster.Cluster, namespaceExclusionList []string, resourceExclusionList []string) {
	e := New(cluster)
	client := e.cluster.Client

	namespaces, err := client.CoreV1().Namespaces().List(context.TODO(), v1.ListOptions{})
	if err != nil || len(namespaces.Items) < 1 {
		return
	}

	for _, namespace := range namespaces.Items {
		for _, resource := range resources {
			e.Watch(namespace.Name, resource)
		}
	}
}
