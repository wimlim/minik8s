package kubectl

import (
	"encoding/json"
	"fmt"
	"io"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/apirequest"
	"minik8s/pkg/config/apiconfig"
	"net/http"
	"strings"

	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Show details of a specific resource or group of resources",
	Run:   describeHandler,
}

func describeHandler(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("no args")
		return
	}

	fd, err := os.Open(args[0])
	if err != nil {
		fmt.Println("open file error")
		return
	}
	defer fd.Close()
	content, err := io.ReadAll(fd)
	if err != nil {
		fmt.Println("read file error")
		return
	}

	kind, err := parseApiObjKind(content)
	if err != nil {
		fmt.Println("parse api obj error")
		return
	}

	var apiObject apiobj.ApiObject
	switch kind {
	case "Node":
		apiObject = &apiobj.Node{}
	case "Pod":
		apiObject = &apiobj.Pod{}
	case "Service":
		describeService(content)
		return
	case "ReplicaSet":
		apiObject = &apiobj.ReplicaSet{}
	case "Hpa":
		apiObject = &apiobj.Hpa{}
	case "Dns":
		apiObject = &apiobj.Dns{}
	case "Function":
		apiObject = &apiobj.Function{}
	case "Workflow":
		apiObject = &apiobj.Workflow{}
	case "PV":
		apiObject = &apiobj.PV{}
	case "PVC":
		apiObject = &apiobj.PVC{}

	}

	describeApiObject(content, apiObject)

}

func describeApiObject(content []byte, apiObject apiobj.ApiObject) {

	err := yaml.Unmarshal(content, apiObject)
	// fmt.Println(string(content))
	if err != nil {
		fmt.Printf("unmarshal %s error\n", apiObject.GetKind())
		return
	}
	if apiObject.GetNamespace() == "" {
		if apiObject.GetKind() != "Node" {
			apiObject.SetNamespace("default")
		}
	}

	URL := apiconfig.Kind2URL[apiObject.GetKind()]
	if apiObject.GetKind() != "Node" {
		URL = strings.Replace(URL, ":namespace", apiObject.GetNamespace(), -1)
	}
	URL = strings.Replace(URL, ":name", apiObject.GetName(), -1)
	HttpUrl := apiconfig.GetApiServerUrl() + URL

	fmt.Println("Get " + HttpUrl)

	response, err := http.Get(HttpUrl)
	if err != nil {
		fmt.Printf("get %s error", apiObject.GetKind())
		return
	}
	defer response.Body.Close()

	var res map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		fmt.Printf("decode %s error\n", apiObject.GetKind())
		return
	}
	data := res["data"].(string)
	err = json.Unmarshal([]byte(data), apiObject)
	if err != nil {
		fmt.Printf("unmarshal %s error\n", apiObject.GetKind())
		return
	}

	podJson, err := json.MarshalIndent(apiObject, "", "    ")
	if err != nil {
		fmt.Printf("marshal %s error\n", apiObject.GetKind())
		return
	}
	fmt.Println(string(podJson))
}

func describeService(content []byte) {
	var service apiobj.Service
	err := yaml.Unmarshal(content, &service)
	if err != nil {
		fmt.Println("unmarshal service error")
		return
	}

	URL := apiconfig.URL_Service
	URL = strings.Replace(URL, ":namespace", "default", -1)
	URL = strings.Replace(URL, ":name", service.MetaData.Name, -1)
	HttpUrl := apiconfig.GetApiServerUrl() + URL

	fmt.Println("Get " + HttpUrl)

	response, err := http.Get(HttpUrl)
	if err != nil {
		fmt.Printf("get service error")
		return
	}
	defer response.Body.Close()

	var res map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		fmt.Printf("decode service error\n")
		return
	}
	data := res["data"].(string)
	err = json.Unmarshal([]byte(data), &service)
	if err != nil {
		fmt.Printf("unmarshal sevice error\n")
		return
	}

	fmt.Printf("name:\n"+"\t%s\n", service.MetaData.Name)
	fmt.Printf("selector:\n"+"\tsvc:%s\n", service.Spec.Selector["svc"])
	fmt.Printf("vip:\n"+"\t%s\n", service.Spec.ClusterIP)
	fmt.Printf("ports(Port/TargetPort/NodePort):\n")
	for _, port := range service.Spec.Ports {
		fmt.Printf("\t%d/%d/%d\n", port.Port, port.TargetPort, port.NodePort)
	}
	fmt.Printf("endpoints:\n")
	pods, _ := apirequest.GetAllPods()
	for _, pod := range pods {
		if pod.MetaData.Labels["svc"] == service.Spec.Selector["svc"] {
			if pod.Status.PodIP != "" {
				fmt.Printf("\t%s\n", pod.Status.PodIP)
			}
		}
	}
}
