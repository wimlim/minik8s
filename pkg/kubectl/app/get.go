package kubectl

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/apirequest"
	"minik8s/pkg/config/apiconfig"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Display one or many resources",
	Run:   getHandler,
}

func getHandler(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("no args")
		return
	}
	kind := args[0]

	if kind == "pods" {
		pods, err := apirequest.GetAllPods()
		if err != nil {
			fmt.Println("Error getting pods:", err)
		}
		for _, pod := range pods {
			jsonData, _ := json.MarshalIndent(pod, "", "    ")
			fmt.Println(string(jsonData))
		}
		return
	}
	if kind == "nodes" {
		nodes, err := apirequest.GetAllNodes()
		if err != nil {
			fmt.Println("Error getting nodes:", err)
		}
		for _, node := range nodes {
			jsonData, _ := json.MarshalIndent(node, "", "    ")
			fmt.Println(string(jsonData))
		}
		return
	}

	var apiObject apiobj.ApiObject
	switch kind {
	case "Pod":
		apiObject = &apiobj.Pod{}
	case "Service":
		apiObject = &apiobj.Service{}
	case "ReplicaSet":
		apiObject = &apiobj.ReplicaSet{}
	}

	getApiObject(args[1], apiObject)
}

func getApiObject(arg string, apiObject apiobj.ApiObject) {
	namespace_obj := strings.Split(arg, "/")
	namespace_name := namespace_obj[0]
	obj_name := namespace_obj[1]

	URL := apiconfig.Kind2URL[apiObject.GetKind()]
	URL = strings.Replace(URL, ":namespace", namespace_name, -1)
	URL = strings.Replace(URL, ":name", obj_name, -1)
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
