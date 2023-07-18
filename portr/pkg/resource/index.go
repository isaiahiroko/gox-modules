package resource

import (
	"fmt"

	"github.com/origine-run/portr/pkg/cluster"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	res "k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

type Item = []byte

type Items = []Item

type resource struct {
	cluster *cluster.Cluster
}

func (r *resource) Get(perPage int) Items {

}

func (r *resource) Apply(item Item) {
	client := r.cluster.Client
	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode(item, nil, nil)
	if err != nil {
		fmt.Printf("%#v", err)
	}

	// Create a REST mapper that tracks information about the available resources in the cluster.
	groupResources, err := restmapper.GetAPIGroupResources(client.Discovery())
	if err != nil {
		return
	}

	rm := restmapper.NewDiscoveryRESTMapper(groupResources)

	// Get some metadata needed to make the REST request.
	gvk := obj.GetObjectKind().GroupVersionKind()
	gk := schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind}
	mapping, err := rm.RESTMapping(gk, gvk.Version)
	if err != nil {
		return
	}

	// name, err := meta.NewAccessor().Name(obj)
	// if err != nil {
	// 	return
	// }

	// Create a client specifically for creating the object.
	gv := mapping.GroupVersionKind.GroupVersion()
	restConfig := &rest.Config{}
	restConfig.ContentConfig = res.UnstructuredPlusDefaultContentConfig()
	restConfig.GroupVersion = &gv
	if len(gv.Group) == 0 {
		restConfig.APIPath = "/api"
	} else {
		restConfig.APIPath = "/apis"
	}

	restClient, err := rest.RESTClientFor(restConfig)
	if err != nil {
		return
	}

	// Use the REST helper to create the object in the "default" namespace.
	restHelper := res.NewHelper(restClient, mapping)
	return restHelper.Apply("default", false, obj, &metaV1.ApplyOptions{})
}
