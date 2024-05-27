package apiserverutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/config/apiconfig"
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
