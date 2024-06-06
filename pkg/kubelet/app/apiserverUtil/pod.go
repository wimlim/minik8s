package apiserverutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/apirequest"
	"minik8s/pkg/config/apiconfig"
	"minik8s/pkg/minik8sTypes"
	"net/http"
	"strings"
)

func PodUpdate(pod *apiobj.Pod) {
	URL := apiconfig.URL_Pod
	URL = strings.Replace(URL, ":namespace", pod.MetaData.Namespace, -1)
	URL = strings.Replace(URL, ":name", pod.MetaData.Name, -1)
	HttpUrl := apiconfig.GetApiServerUrl() + URL
	jsonData, err := json.Marshal(pod)
	if err != nil {
		fmt.Println("marshal pod error")
		return
	}
	req, err := http.NewRequest(http.MethodPut, HttpUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("create put request error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("put error:", err)
		return
	}
	defer response.Body.Close()
}

func PodStatusUpdate(podIdentifier minik8sTypes.PodIdentifier, podStatus *apiobj.PodStatus) {
	URL := apiconfig.URL_PodStatus
	URL = strings.Replace(URL, ":namespace", podIdentifier.PodNamespace, -1)
	URL = strings.Replace(URL, ":name", podIdentifier.PodName, -1)
	HttpUrl := apiconfig.GetApiServerUrl() + URL
	jsonData, err := json.Marshal(podStatus)
	if err != nil {
		fmt.Println("marshal pod error")
		return
	}
	req, err := http.NewRequest(http.MethodPut, HttpUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("create put request error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("put error:", err)
		return
	}
	defer response.Body.Close()
}

func GetAllRemotePods() ([]apiobj.Pod, error) {
	return apirequest.GetAllPods()
}

func NodeStatusUpdate(nodeStatus apiobj.NodeStatus, hostNode *apiobj.Node) {
	URL := apiconfig.URL_NodeStatus
	URL = strings.Replace(URL, ":namespace", hostNode.GetNamespace(), -1)
	URL = strings.Replace(URL, ":name", hostNode.GetName(), -1)
	HttpUrl := apiconfig.GetApiServerUrl() + URL
	jsonData, err := json.Marshal(nodeStatus)
	if err != nil {
		fmt.Println("marshal node error")
		return
	}
	req, err := http.NewRequest(http.MethodPut, HttpUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("create put request error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("put error:", err)
		return
	}
	defer response.Body.Close()
}
