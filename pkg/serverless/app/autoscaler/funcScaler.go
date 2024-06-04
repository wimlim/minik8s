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

const (
	func_scale_time = time.Second * 30
)

type Record struct {
	func_namespace string
	func_name      string
	call_frequency int
	start_time     time.Time
	end_time       time.Time
}

type FuncScaler struct {
	funcMap   map[string]apiobj.Function
	recordMap map[string]Record
}

func NewFuncScaler() *FuncScaler {
	return &FuncScaler{
		funcMap:   make(map[string]apiobj.Function),
		recordMap: make(map[string]Record),
	}
}

func (fs *FuncScaler) Run() {
	runner.NewRunner().RunLoop(5*time.Second, 5*time.Second, fs.func_routine)
}

func (fs *FuncScaler) func_routine() {

	funcs, err := apirequest.GetAllFunctions()
	if err != nil {
		fmt.Println("get all functions error")
		return
	}

	//get remote func
	remoteFuncMap := make(map[string]bool)
	for _, f := range funcs {
		key := fmt.Sprintf("%s/%s", f.MetaData.Namespace, f.MetaData.Name)
		remoteFuncMap[key] = true
	}

	//delete func that not exist in remote
	for id, f := range fs.funcMap {
		key := fmt.Sprintf("%s/%s", f.MetaData.Namespace, f.MetaData.Name)
		if _, ok := remoteFuncMap[key]; !ok {
			fs.Deletefunc(f)
			delete(fs.funcMap, id)
			key := fmt.Sprintf("%s/%s", f.MetaData.Namespace, f.MetaData.Name)
			delete(fs.recordMap, key)
		}
	}

	//add new func
	for _, f := range funcs {
		key := fmt.Sprintf("%s/%s", f.MetaData.Namespace, f.MetaData.Name)
		if _, ok := fs.funcMap[key]; !ok {
			fs.funcMap[key] = f
			fs.Addfunc(f)
		} else {
			if fs.recordMap[key] != (Record{}) {
				end_time := fs.recordMap[key].end_time
				if time.Now().After(end_time) {

					if fs.recordMap[key].call_frequency == 0 {
						fmt.Println("scale down function")
						fs.DeleteRelica(f)
					} else if fs.recordMap[key].call_frequency > 100 {
						expectSize := fs.recordMap[key].call_frequency/100 + 1
						fmt.Println("scale up function")
						fs.AddReplica(f.MetaData.Namespace, f.MetaData.Name, expectSize)
					} else {
						fs.AddReplica(f.MetaData.Namespace, f.MetaData.Name, 1)
					}

					record := fs.recordMap[key]
					record.start_time = time.Now()
					record.end_time = time.Now().Add(func_scale_time)
					record.call_frequency = 0
					fs.recordMap[key] = record

				}
			}

		}
	}

}

func (fs *FuncScaler) AddRecord(func_namespace string, func_name string) {
	//TODO
	new_call_record := Record{
		func_namespace: func_namespace,
		func_name:      func_name,
		start_time:     time.Now(),
		end_time:       time.Now().Add(func_scale_time),
		call_frequency: 1,
	}

	key := fmt.Sprintf("%s/%s", func_namespace, func_name)
	if _, ok := fs.recordMap[key]; !ok {
		fs.recordMap[key] = new_call_record
	} else {
		record := fs.recordMap[key]
		record.call_frequency++
		fs.recordMap[key] = record
		fmt.Println("call frequency: ", record.call_frequency)
	}

}

func (fs *FuncScaler) AddReplica(func_namespace string, func_name string, num int) {
	obj, _ := apirequest.GetRequest(func_namespace, func_name+"-replica", "ReplicaSet")
	replica := obj.(*apiobj.ReplicaSet)
	replica.Spec.Replicas = num

	URL := apiconfig.GetApiServerUrl() + apiconfig.URL_ReplicaSet
	URL = strings.Replace(URL, ":namespace", replica.MetaData.Namespace, -1)
	URL = strings.Replace(URL, ":name", replica.MetaData.Name, -1)

	err := apirequest.PutRequest(URL, replica)
	if err != nil {
		fmt.Println("put replica error")
		return
	}
}

func (fs *FuncScaler) DeleteRelica(f apiobj.Function) {

	obj, _ := apirequest.GetRequest(f.MetaData.Namespace, f.MetaData.Name+"-replica", "ReplicaSet")
	replica := obj.(*apiobj.ReplicaSet)
	replica.Spec.Replicas = 0

	URL := apiconfig.GetApiServerUrl() + apiconfig.URL_ReplicaSet
	URL = strings.Replace(URL, ":namespace", replica.MetaData.Namespace, -1)
	URL = strings.Replace(URL, ":name", replica.MetaData.Name, -1)

	err := apirequest.PutRequest(URL, replica)
	if err != nil {
		fmt.Println("put replica error")
		return
	}
}

func (fs *FuncScaler) Addfunc(f apiobj.Function) {

	r := registry.NewRegistry()
	if r == nil {
		fmt.Println("NewRegistry error")
	}

	r.BuildImage(f)

	imageName := fmt.Sprintf("func/%s:latest", f.MetaData.Name)
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
			Replicas: 0,
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
						"func_uid":       f.MetaData.UID,
						"func_namespace": f.MetaData.Namespace,
						"func_name":      f.MetaData.Name,
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
func (fs *FuncScaler) Deletefunc(f apiobj.Function) {
	delete(fs.funcMap, f.MetaData.UID)

	URL := apiconfig.GetApiServerUrl() + apiconfig.URL_ReplicaSet
	URL = strings.Replace(URL, ":namespace", f.MetaData.Namespace, -1)
	URL = strings.Replace(URL, ":name", f.MetaData.Name+"-replica", -1)

	apirequest.DeleteRequest(URL)
}
