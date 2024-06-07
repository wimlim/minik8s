package apiserverutil

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"testing"
)

func TestGetAllRemotePods(t *testing.T) {
	allRemotePods, err := GetAllRemotePods()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	remotePodsMap := map[string]*apiobj.Pod{}
	for _, pod := range allRemotePods {
		podId := pod.MetaData.UID
		remotePodsMap[podId] = &pod
	}
	targetPod := remotePodsMap["f25e6e79-940f-43"]
	targetPodJSON, err := json.MarshalIndent(*targetPod, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 打印 JSON 格式的 targetPod
	fmt.Println(string(targetPodJSON))
}
