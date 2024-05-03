package cmd

import (
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/config/apiconfig"
	"net/http"
	"strings"
	"encoding/json"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use: "get",
	Short: "Display one or many resources",
	Run: getHandler,
}

func getHandler(cmd *cobra.Command, args []string){
	if(len(args) == 0){
		fmt.Println("no args")
		return
	}
	kind := args[0]

	switch kind {
	case "Pod":
		getPod(args[1])
	case "Service":
		fmt.Println("get service")
	}
	
}

func getPod(arg string){
	namespace_pod := strings.Split(arg, "/")
	namespace_name := namespace_pod[0]
	pod_name := namespace_pod[1]

	URL := apiconfig.URL_Pod
	URL = strings.Replace(URL,":namespace",namespace_name,-1)
	URL = strings.Replace(URL,":name",pod_name,-1)
	HttpUrl := apiconfig.GetApiServerUrl() + URL

	fmt.Println("Get " + HttpUrl)

	var pod apiobj.Pod
	
	response, err := http.Get(HttpUrl)
	if err != nil {
		fmt.Println("get pod error")
		return
	}
	defer response.Body.Close()
    
	var res map[string] interface{}
	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		fmt.Println("decode pod error")
		return
	}
	data := res["data"].(string)
	
	err = json.Unmarshal([]byte(data), &pod)
	if err != nil {
		fmt.Println("unmarshal pod error")
		return
	}

	fmt.Println(data)
}