package autoscaler

import (
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/apirequest"
	"minik8s/pkg/config/apiconfig"
	"minik8s/pkg/config/serverlessconfig"
	"minik8s/pkg/serverless/app/registry"
	"minik8s/tools/runner"
	"strings"
	"time"
)

type funcScaler struct {
	funcMap map[string]apiobj.Function
}

func NewFuncScaler() *funcScaler {
	return &funcScaler{
		funcMap: make(map[string]apiobj.Function),
	}
}

func (fs *funcScaler) Run() {
	runner.NewRunner().RunLoop(5*time.Second, 5*time.Second, fs.func_routine)
}

func (fs *funcScaler) func_routine() {

	funcs, err := apirequest.GetAllFunctions()
	if err != nil {
		fmt.Println("get all functions error")
		return
	}

	remoteFuncMap := make(map[string]bool)
	for _, f := range funcs {
		remoteFuncMap[f.MetaData.UID] = true
	}

	for id, f := range fs.funcMap {
		if _, ok := remoteFuncMap[f.MetaData.UID]; !ok {
			fs.Deletefunc(f)
			delete(fs.funcMap, id)
		}
	}

	for _, f := range funcs {
		if _, ok := fs.funcMap[f.MetaData.UID]; !ok {
			fs.funcMap[f.MetaData.UID] = f
			fs.Addfunc(f)
		}
	}

}

func (fs *funcScaler) Addfunc(f apiobj.Function) {

	r := registry.NewRegistry()
	if r == nil {
		fmt.Println("NewRegistry error")
	}

	r.BuildImage(f)

	imageName := fmt.Sprintf("func/%s:latest", f.MetaData.UID)
	imageRef := fmt.Sprintf("%s/%s", serverlessconfig.GetRegistryServerUrl(), imageName)

	replica := &apiobj.ReplicaSet{
		ApiVersion: "v1",
		Kind:       "Replica",
		MetaData: apiobj.MetaData{
			Name:      f.MetaData.Name + "-replica",
			Namespace: f.MetaData.Namespace,
			Labels: map[string]string{
				"func_uid": f.MetaData.UID,
			},
		},
		Spec: apiobj.ReplicaSetSpec{
			Replicas: 2,
			Selector: apiobj.ReplicaSetSelector{
				MatchLabels: map[string]string{
					"func_uid": f.MetaData.UID,
				},
			},
			Template: apiobj.PodTemplateSpec{
				MetaData: apiobj.MetaData{
					Name:      f.MetaData.Name,
					Namespace: f.MetaData.Namespace,
					Labels: map[string]string{
						"func_uid": f.MetaData.UID,
						"func_namespace": f.MetaData.Namespace,
						"func_name": f.MetaData.Name,
					},
				},
				Spec: apiobj.PodSpec{
					Containers: []apiobj.Container{
						{
							Name:  f.MetaData.Name,
							Image: imageRef,
							Ports: []apiobj.ContainerPort{
								{
									ContainerPort: 8080,
								},
							},
						},
					},
				},
			},
		},
	}

	URL := apiconfig.GetApiServerUrl() + apiconfig.URL_ReplicaSet
	URL = strings.Replace(URL, ":namespace", replica.MetaData.Namespace, -1)
	URL = strings.Replace(URL, ":name", replica.MetaData.Name, -1)

	err := apirequest.PostRequest(URL, replica)
	if err != nil {
		fmt.Println("post replica error")
		return
	}

}
func (fs *funcScaler) Deletefunc(f apiobj.Function) {
	delete(fs.funcMap, f.MetaData.UID)

	URL := apiconfig.GetApiServerUrl() + apiconfig.URL_ReplicaSet
	URL = strings.Replace(URL, ":namespace", f.MetaData.Namespace, -1)
	URL = strings.Replace(URL, ":name", f.MetaData.Name+"-replica", -1)

	apirequest.DeleteRequest(URL)
}
