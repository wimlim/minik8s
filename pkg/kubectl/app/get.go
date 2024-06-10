package kubectl

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/apirequest"
	"minik8s/pkg/config/apiconfig"
	"net/http"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
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
		PrintPodsTable(pods)
		return
	}
	if kind == "nodes" {
		nodes, err := apirequest.GetAllNodes()
		if err != nil {
			fmt.Println("Error getting nodes:", err)
		}
		PrintNodesTable(nodes)
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
	case "Job":
		apiObject = &apiobj.Job{}
	case "Result":
		getJobResult(args[1])
		return
	}

	getApiObject(args[1], apiObject, kind)
}

func getApiObject(arg string, apiObject apiobj.ApiObject, kind string) {
	namespace_obj := strings.Split(arg, "/")
	namespace_name := namespace_obj[0]
	obj_name := namespace_obj[1]

	URL := apiconfig.Kind2URL[kind]
	URL = strings.Replace(URL, ":namespace", namespace_name, -1)
	URL = strings.Replace(URL, ":name", obj_name, -1)
	HttpUrl := apiconfig.GetApiServerUrl() + URL

	fmt.Println("Get " + HttpUrl)

	response, err := http.Get(HttpUrl)
	if err != nil {
		fmt.Printf("get %s error", kind)
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

func PrintNodesTable(nodes []apiobj.Node) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Node Name", "ApiVersion", "Address"})

	for _, node := range nodes {
		t.AppendRow(table.Row{
			node.MetaData.Name,
			node.ApiVersion,
			node.IP,
		})
	}
	t.Render()
}
func PrintPodsTable(pods []apiobj.Pod) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Pod Name", "Node", "Ip", "Containers"})

	// 添加数据行
	for _, pod := range pods {
		containers := ""
		for _, container := range pod.Spec.Containers {
			containers += fmt.Sprintf("%s (%s)\n", container.Name, container.Image)
		}
		// 删除最后一个换行符
		containers = strings.TrimSuffix(containers, "\n")

		t.AppendRow([]interface{}{pod.MetaData.Name, pod.Spec.NodeName, pod.Status.PodIP, containers})
	}

	t.Render()
}

func getJobResult(name string){

	namespace := "default"

	URL := apiconfig.URL_Job
	URL = strings.Replace(URL, ":namespace", namespace, -1)
	URL = strings.Replace(URL, ":name", name, -1)
	HttpUrl := apiconfig.GetApiServerUrl() + URL

	response, err := http.Get(HttpUrl)
	if err != nil {
		fmt.Printf("get result error")
		return
	}
	defer response.Body.Close()

	var res map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		fmt.Printf("decode result error\n")
		return
	}
	data := res["data"].(string)

	var job apiobj.Job
	err = json.Unmarshal([]byte(data), &job)
	if err != nil {
		fmt.Printf("unmarshal result error\n")
		return
	}

	if(job.Status.Phase == "Running"){
		fmt.Println("Job is running, please wait for a moment.")
		return
	}else{
		fmt.Println(job.Status.Result)
	}

}