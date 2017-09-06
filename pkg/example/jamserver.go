package example

import (
	"encoding/json"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/meta"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/pkg/runtime"
	"k8s.io/client-go/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
)

type JamServerSpec struct {
	Owner   string `json:"owner"`
	Message string `json:"message"`
}

type JamServer struct {
	unversioned.TypeMeta    `json:",inline"`
	Metadata api.ObjectMeta `json:"metadata"`

	Spec JamServerSpec `json:"spec"`
}

type JamServerList struct {
	unversioned.TypeMeta          `json:",inline"`
	Metadata unversioned.ListMeta `json:"metadata"`

	Items []JamServer `json:"items"`
}

func NewJamServerClientForConfig(c *rest.Config) (*rest.RESTClient, error) {
	c.APIPath = "/apis"
	c.GroupVersion = &unversioned.GroupVersion{
		Group:   "gojam.server.op",
		Version: "v1",
	}
	c.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: api.Codecs}

	schemeBuilder := runtime.NewSchemeBuilder(
		func(scheme *runtime.Scheme) error {
			scheme.AddKnownTypes(
				*c.GroupVersion,
				&JamServer{},
				&JamServerList{},
				&api.ListOptions{},
				&api.DeleteOptions{},
			)
			return nil
		})
	schemeBuilder.AddToScheme(api.Scheme)

	return rest.RESTClientFor(c)
}

func EnsureJamServerThirdPartyResource(client *kubernetes.Clientset) error {
	_, err := client.Extensions().ThirdPartyResources().Get("jamservers.gojam.server.op")
	if err != nil {
		// The resource doesn't exist, so we create it.
		newResource := v1beta1.ThirdPartyResource{
			ObjectMeta: v1.ObjectMeta{
				Name: "jamservers.gojam.server.op",
			},
			Description: "A specification of a conversation",
			Versions: []v1beta1.APIVersion{
				{Name: "v1"},
			},
		}

		_, err = client.Extensions().ThirdPartyResources().Create(&newResource)
	}

	return err
}

// The code below is used only to work around a known problem with third-party
// resources and ugorji. If/when these issues are resolved, the code below
// should no longer be required.
//

func (c *JamServer) GetObjectKind() unversioned.ObjectKind {
	return &c.TypeMeta
}

func (c *JamServer) GetObjectMeta() meta.Object {
	return &c.Metadata
}

func (cl *JamServerList) GetObjectKind() unversioned.ObjectKind {
	return &cl.TypeMeta
}

func (cl *JamServerList) GetListMeta() unversioned.List {
	return &cl.Metadata
}

type JamServerListCopy JamServerList
type JamServerCopy JamServer

func (e *JamServer) UnmarshalJSON(data []byte) error {
	tmp := JamServerCopy{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	tmp2 := JamServer(tmp)
	*e = tmp2
	return nil
}

func (el *JamServerList) UnmarshalJSON(data []byte) error {
	tmp := JamServerListCopy{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	tmp2 := JamServerList(tmp)
	*el = tmp2
	return nil
}
