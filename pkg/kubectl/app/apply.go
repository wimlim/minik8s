package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/config/apiconfig"
	"net/http"
	"os"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"strings"
)

var applyCmd = &cobra.Command{
	Use: "apply",
	Short: "Apply a configuration to a resource by filename or stdin",
	Run:	applyHandler,	
}

//kubectl apply pod.yaml
func applyHandler(cmd *cobra.Command, args []string){
	if(len(args) == 0){
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
		applyPod(content)
	case "Service":
		fmt.Println("apply service")
	}
	
}
func applyPod(content []byte){
	var pod apiobj.Pod
	err := yaml.Unmarshal(content, &pod)
	if err != nil {
		fmt.Println("unmarshal pod error")
		return
	}
	if(pod.MetaData.Namespace == ""){
		pod.MetaData.Namespace = "default"
	}
	
	URL := apiconfig.URL_Pod
	URL = strings.Replace(URL,":namespace",pod.MetaData.Namespace,-1)
	URL = strings.Replace(URL,":name",pod.MetaData.Name,-1)
	HttpUrl := apiconfig.GetApiServerUrl() + URL
	fmt.Println("Post " +  HttpUrl)
	jsonData, err := json.Marshal(pod)
	//fmt.Println(string(jsonData))
	if err != nil {
		fmt.Println("marshal pod error")
		return
	}
	response, err := http.Post(HttpUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("post error")
		return
	}
	defer response.Body.Close()
	
}