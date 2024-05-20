package main

import (
	"minik8s/pkg/apirequest"
	"minik8s/pkg/apiobj"
	"encoding/json"
	"fmt"
)

func main() {
	// var pod apiobj.Pod

	// apirequest.GetApiObjectAfterRequest("/api/v1/namespaces/:namespace/pods/:name",
	// 	&pod)

	// jsonData, _ := json.Marshal(pod)
	// fmt.Println(string(jsonData))

	// pods, err := apirequest.GetAllPods()
	// if err != nil {
	// 	fmt.Println("Error getting pods:", err)
	// 	return
	// }
	// for _, pod := range pods {
	// 	jsonData, _ := json.Marshal(pod)
	// 	fmt.Println(string(jsonData))
	// }

	// for i := 0; i < 10; i++ {
	// 	var ip = serviceconfig.AllocateIp()
	// 	if ip != "" {
	// 		fmt.Println(ip)
	// 	}
	// }
	
	var pod apiobj.ApiObject
	pod, err := apirequest.GetRequest("http://127.0.0.1:8080/api/v1/namespaces/default/pods/http-server", "Pod")
	pod = pod.(*apiobj.Pod)
	if err != nil {
		return
	}
	jsonData, _ := json.Marshal(pod)
	fmt.Println(string(jsonData))
}
