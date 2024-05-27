package kubectl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/config/apiconfig"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a configuration to a resource by filename or stdin",
	Run:   applyHandler,
}

// kubectl apply pod.yaml
func applyHandler(cmd *cobra.Command, args []string) {
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
		apiObject = &apiobj.Service{}
	case "ReplicaSet":
		apiObject = &apiobj.ReplicaSet{}
	case "Hpa":
		apiObject = &apiobj.Hpa{}
	case "Dns":
		apiObject = &apiobj.Dns{}
	case "Function":
		applyFunction(content)
		return
	case "Workflow":
		apiObject = &apiobj.Workflow{}
	case "PV":
		apiObject = &apiobj.PV{}
	case "PVC":
		apiObject = &apiobj.PVC{}
	}

	applyApiObject(content, apiObject)
}
func applyApiObject(content []byte, apiObject apiobj.ApiObject) {

	err := yaml.Unmarshal(content, apiObject)
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
	fmt.Println("Post " + HttpUrl)
	jsonData, err := json.Marshal(apiObject)
	// fmt.Println(string(jsonData))
	if err != nil {
		fmt.Printf("marshal %s error\n", apiObject.GetKind())
		return
	}
	response, err := http.Post(HttpUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("post error")
		return
	}
	defer response.Body.Close()
	fmt.Printf("apply %s request sent\n", apiObject.GetKind())
}

func applyFunction(content []byte) {
	var function apiobj.Function
	err := yaml.Unmarshal(content, &function)
	if err != nil {
		fmt.Println("unmarshal function error")
		return
	}
	if function.GetNamespace() == "" {
		function.SetNamespace("default")
	}

	filepath := function.Spec.Path
	fd, err := os.Open(filepath)
	if err != nil {
		fmt.Println("open function file error")
		return
	}
	defer fd.Close()

	content, err = io.ReadAll(fd)
	if err != nil {
		fmt.Println("read function file error")
		return
	}

	function.Spec.Content = content
	fmt.Println(string(content))

	URL := apiconfig.URL_Function
	URL = strings.Replace(URL, ":namespace", function.GetNamespace(), -1)
	URL = strings.Replace(URL, ":name", function.GetName(), -1)
	HttpUrl := apiconfig.GetApiServerUrl() + URL
	fmt.Println("Post " + HttpUrl)
	jsonData, err := json.Marshal(function)
	if err != nil {
		fmt.Println("marshal function error")
		return
	}
	response, err := http.Post(HttpUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("post error")
		return
	}
	defer response.Body.Close()
	fmt.Println("apply function request sent")
}
