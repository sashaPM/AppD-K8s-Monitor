package watchers

import (
	"fmt"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"

	"github.com/sjeltuhin/clusterAgent/utils"
)

type EndpointWatcher struct {
	Client        *kubernetes.Clientset
	EndpointCache map[string]v1.Endpoints
}

func NewEndpointWatcher(client *kubernetes.Clientset) *EndpointWatcher {
	epw := EndpointWatcher{Client: client, EndpointCache: make(map[string]v1.Endpoints)}
	return &epw
}

//end points
func (pw EndpointWatcher) WatchEndpoints() {
	api := pw.Client.CoreV1()
	listOptions := metav1.ListOptions{}
	fmt.Println("Starting Endpoint Watcher...")

	watcher, err := api.Endpoints(metav1.NamespaceAll).Watch(listOptions)
	if err != nil {
		fmt.Printf("Issues when setting up endpoint watcher. %v", err)
	}

	ch := watcher.ResultChan()

	for ev := range ch {
		ep, ok := ev.Object.(*v1.Endpoints)
		if !ok {
			fmt.Printf("Expected endpoints, but received an object of an unknown type. ")
			continue
		}
		switch ev.Type {
		case watch.Added:
			pw.onNewEndpoint(ep)
			break

		case watch.Deleted:
			pw.onDeleteEndpoint(ep)
			break

		case watch.Modified:
			pw.onUpdateEndpoint(ep)
			break
		}

	}
	fmt.Println("Exiting endpoint watcher.")
}

func (pw EndpointWatcher) onNewEndpoint(ep *v1.Endpoints) {
	pw.EndpointCache[utils.GetEndpointKey(ep)] = *ep
}

func (pw EndpointWatcher) onDeleteEndpoint(ep *v1.Endpoints) {
	_, ok := pw.EndpointCache[utils.GetEndpointKey(ep)]
	if ok {
		delete(pw.EndpointCache, utils.GetEndpointKey(ep))
	}
}

func (pw EndpointWatcher) onUpdateEndpoint(ep *v1.Endpoints) {
	pw.EndpointCache[utils.GetEndpointKey(ep)] = *ep
}