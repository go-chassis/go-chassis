package kuberegistry

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/stretchr/testify/assert"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"
)

func TestDiscoveryController(t *testing.T) {
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	client := fake.NewSimpleClientset()
	sharedfactory := informers.NewSharedInformerFactory(client, 0)
	sInformer := sharedfactory.Core().V1().Services()
	eInformer := sharedfactory.Core().V1().Endpoints()

	dc := NewDiscoveryController(sInformer, eInformer, client)
	sharedfactory.Start(ctx.Done())
	dc.Run(ctx.Done())

	// create endpoints
	p := &v1.Endpoints{ObjectMeta: metav1.ObjectMeta{Name: "kubeserver"},
		Subsets: []v1.EndpointSubset{{
			Addresses: []v1.EndpointAddress{{IP: "127.0.0.1",
				TargetRef: &v1.ObjectReference{UID: "12345"}}},
			Ports: []v1.EndpointPort{{Name: "rest", Port: 9090}},
		}}}
	_, err := client.Core().Endpoints("default").Create(p)
	if err != nil {
		t.Errorf("error create endpoints: %v", err)
	}

	// create services
	s := &v1.Service{ObjectMeta: metav1.ObjectMeta{Name: "kubeserver"}}
	_, err = client.Core().Services("default").Create(s)
	if err != nil {
		t.Errorf("error create service: %v", err)
	}

	time.Sleep(1 * time.Second)
	ret, err := dc.FindEndpoints("kubeserver.default.svc.local", nil)
	assert.Equal(t, len(ret), 1)
	for _, ep := range ret {
		log.Printf("Got endpoints %s(%s)", ep.EndpointsMap, ep.ServiceID)
		assert.Equal(t, ep.EndpointsMap["rest"], "127.0.0.1:9090")
		assert.Equal(t, ep.ServiceID, "kubeserver.default")
	}

	svc, err := dc.GetAllServices()
	assert.Equal(t, len(svc), 1)
	for _, ss := range svc {
		assert.Equal(t, ss.ServiceName, "kubeserver")
	}
}
