package apirequest

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/config/apiconfig"
	"net/http"
)

func GetAllPods() ([]apiobj.Pod, error) {
	URL := apiconfig.URL_AllPods
	HttpURL := apiconfig.GetApiServerUrl() + URL

	response, err := http.Get(HttpURL)
	if err != nil {
		fmt.Println("HTTP request error:", err)
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("HTTP request returned status code:", response.StatusCode)
		return nil, fmt.Errorf("status code: %d", response.StatusCode)
	}

	var res map[apiobj.Pod]interface{}
	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		fmt.Println("decode pod error")
		return nil, fmt.Errorf("status code: %d", response.StatusCode)
	}

	data, ok := res["data"].([]interface{})
	if !ok {
		fmt.Println("expected type []interface{} for field 'data', got something else")
		return nil, fmt.Errorf("type assertion failed for 'data'")
	}

	var pods []apiobj.Pod
	for _, item := range data {
		pod, ok := item.(apiobj.Pod)
		if !ok {
			fmt.Println("type assertion failed for an item in 'data'")
			return nil, fmt.Errorf("type assertion failed for an item in 'data'")
		}
		pods = append(pods, pod)
	}
	return pods, nil
}

func GetAllNodes() ([]string, error) {
	URL := apiconfig.URL_AllNodes
	HttpURL := apiconfig.GetApiServerUrl() + URL

	response, err := http.Get(HttpURL)
	if err != nil {
		fmt.Println("HTTP request error:", err)
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("HTTP request returned status code:", response.StatusCode)
		return nil, fmt.Errorf("status code: %d", response.StatusCode)
	}

	var res map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		fmt.Println("decode pod error")
		return nil, err
	}

	data, ok := res["data"].([]interface{})
	if !ok {
		fmt.Println("expected type []interface{} for field 'data', got something else")
		return nil, fmt.Errorf("type assertion failed for 'data'")
	}

	// 将 interface{} 列表转换为字符串列表
	var nodes []string
	for _, item := range data {
		str, ok := item.(string)
		if !ok {
			fmt.Println("type assertion failed for an item in 'data'")
			return nil, fmt.Errorf("type assertion failed for an item in 'data'")
		}
		nodes = append(nodes, str)
	}

	return nodes, nil
}
