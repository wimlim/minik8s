package kubectl

import (
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

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete resources by filenames, stdin, resources and names",
	Run:   deleteHandler,
}

// kubectl delete pod.yaml
func deleteHandler(cmd *cobra.Command, args []string) {
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
	}

	deleteApiObject(content, apiObject)
}
func deleteApiObject(content []byte, apiObject apiobj.ApiObject) {
	err := yaml.Unmarshal(content, apiObject)
	if err != nil {
		fmt.Printf("unmarshal %s error\n",apiObject.GetKind())
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

	fmt.Println("Delete " + HttpUrl)
	request, err := http.NewRequest("DELETE", HttpUrl, nil)
	if err != nil {
		fmt.Println("new request error")
		return
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("do request error")
		return
	}
	defer resp.Body.Close()
	fmt.Printf("delete %s request sent\n", apiObject.GetKind())
}
