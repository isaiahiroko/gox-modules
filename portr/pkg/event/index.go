package event

import (
	"context"

	"github.com/origine-run/portr/pkg/cluster"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type event struct {
	cluster cluster.Cluster
}

func (e *event) Watch(namespace string, resource string) {
	client := e.cluster.Client

	var watcher watch.Interface
	var err error

	switch resource {
	case "service":
		watcher, err = client.CoreV1().Services(namespace).Watch(context.TODO(), metaV1.ListOptions{})
	default:
		return
	}

	if err != nil {
		return
	}

	for event := range watcher.ResultChan() {
		// svc := event.Object.(*v1.Service)

		// switch event.Type {
		// case watch.Added:
		// 	fmt.Printf("Service %s/%s added", svc.ObjectMeta.Namespace, svc.ObjectMeta.Name)
		// case watch.Modified:
		// 	fmt.Printf("Service %s/%s modified", svc.ObjectMeta.Namespace, svc.ObjectMeta.Name)
		// case watch.Deleted:
		// 	fmt.Printf("Service %s/%s deleted", svc.ObjectMeta.Namespace, svc.ObjectMeta.Name)
		// }
	}
}
