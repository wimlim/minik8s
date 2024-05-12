package main

import (
	"fmt"
	"minik8s/pkg/config/serviceconfig"
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

	for i := 0; i < 10; i++ {
		var ip = serviceconfig.AllocateIp()
		if ip != "" {
			fmt.Println(ip)
		}
	}

}
