package kubehandler

import (
	"testing"

	routev1 "github.com/openshift/api/route/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestShouldGetResourceVersionFromK8sObjects(t *testing.T) {
	var route interface{}
	route = routev1.Route{ObjectMeta: metav1.ObjectMeta{ResourceVersion: "test"}}
	assert.NotNil(t, route)

	version, ok := resourceVersion(route)
	assert.True(t, ok)
	assert.Equal(t, "test", version)
}
