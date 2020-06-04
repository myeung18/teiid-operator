package openshift

import (
	"context"
	"fmt"
	"github.com/RHsyseng/operator-utils/pkg/test"
	consolev1 "github.com/openshift/api/console/v1"
	routev1 "github.com/openshift/api/route/v1"
	"github.com/stretchr/testify/assert"
	"github.com/teiid/teiid-operator/pkg/apis/teiid/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

func getInstance(nsName types.NamespacedName) *v1alpha1.VirtualDatabase {
	return &v1alpha1.VirtualDatabase{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nsName.Name,
			Namespace: nsName.Namespace,
		},
	}
}

func TestTest(t *testing.T) {
	apiObjects := []runtime.Object{&v1alpha1.VirtualDatabase{}, &v1alpha1.VirtualDatabaseList{}}

	v1alpha1.SchemeBuilder.Register(apiObjects...)

	scheme, _ := v1alpha1.SchemeBuilder.Build()
	client := fake.NewFakeClientWithScheme(scheme)
	log.Debugf("Fake client created as %v", client)

}

func TestCreateMockServiceClient(t *testing.T) {
	crNamespace := types.NamespacedName{
		Name:      "test",
		Namespace: "vdb-ns",
	}
	vdb := getInstance(crNamespace)
	_ = vdb

	service := MockService()
	err := service.Create(context.TODO(), vdb)
	assert.Nil(t, err)


	bcRoute := &routev1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Name :"",
			Namespace: "",
		},
		Spec: routev1.RouteSpec{
			Host:              "www.example.com",
			Subdomain:         "",
			Path:              "",
			To:                routev1.RouteTargetReference{},
			WildcardPolicy:    "",
		},
	}
	CreateConsoleLink(context.TODO(), bcRoute, service.Client, vdb)

	consoleLinkName := fmt.Sprintf("%s-%s", vdb.ObjectMeta.Name, vdb.Namespace)
	consoleLink := &consolev1.ConsoleLink{}
	err = service.Get(context.TODO(), types.NamespacedName{Name: consoleLinkName}, consoleLink)

	fmt.Println(consoleLink)
	assert.Nil(t, err)

}

func TestConsoleLinkCreation(t *testing.T) {
	crNamespace := types.NamespacedName{
		Name:      "testns",
		Namespace: "cr",
	}
	vdb := getInstance(crNamespace)
	_ = vdb
	//fmt.Println(cr)

	var localSchemeBuilder = runtime.SchemeBuilder{}
	st := test.NewMockPlatformServiceBuilder(localSchemeBuilder)

	apiObjects := []runtime.Object{&v1alpha1.VirtualDatabase{}, &v1alpha1.VirtualDatabaseList{}}
	st.WithScheme(apiObjects...)
	service := st.Build()

	er := service.Create(context.TODO(), vdb)
	fmt.Println(er)

	bcRoute := &routev1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Name :"",
			Namespace: "",
		},
		Spec: routev1.RouteSpec{
			Host:              "www.example.com",
			Subdomain:         "",
			Path:              "",
			To:                routev1.RouteTargetReference{},
			WildcardPolicy:    "",
		},
	}
	CreateConsoleLink(context.TODO(), bcRoute, service.Client, vdb)

	consoleLinkName := fmt.Sprintf("%s-%s", vdb.ObjectMeta.Name, vdb.Namespace)
	consoleLink := &consolev1.ConsoleLink{}
	err := service.Get(context.TODO(), types.NamespacedName{Name: consoleLinkName}, consoleLink)
	assert.Nil(t, err)
	fmt.Println(consoleLink.Spec)
	_ = service


}
