package cmd

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
	switch kind {
	case "Pod":
		deletePod(content)
	case "service":
		fmt.Println("delete service")
	}
}
func deletePod(content []byte) {
	var pod apiobj.Pod
	err := yaml.Unmarshal(content, &pod)
	if err != nil {
		fmt.Println("unmarshal pod error")
		return
	}
	URL := apiconfig.URL_Pod
	if pod.MetaData.Namespace == "" {
		pod.MetaData.Namespace = "default"
	}
	URL = strings.Replace(URL, ":namespace", pod.MetaData.Namespace, -1)
	URL = strings.Replace(URL, ":name", pod.MetaData.Name, -1)
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
	fmt.Println("delete pod success")
}
