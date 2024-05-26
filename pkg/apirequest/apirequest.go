package apirequest

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/config/apiconfig"
	"net/http"
)

func GetAllPods() ([]apiobj.Pod, error) {
	URL := apiconfig.URL_GlobalPods
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
		return nil, fmt.Errorf("status code: %d", response.StatusCode)
	}

	if res["data"] == nil {
		fmt.Println("empty pod list")
		return []apiobj.Pod{}, nil
	}

	data, ok := res["data"].([]interface{})
	if !ok {
		fmt.Println("expected type []interface{} for field 'data', got something else")
		return nil, fmt.Errorf("type assertion failed for 'data'")
	}

	var pods []apiobj.Pod
	for _, item := range data {
		podStr, ok := item.(string)
		if !ok {
			fmt.Println("type assertion failed for an item in 'data'")
			return nil, fmt.Errorf("type assertion failed for an item in 'data'")
		}
		var pod apiobj.Pod
		json.Unmarshal([]byte(podStr), &pod)
		pods = append(pods, pod)
	}
	return pods, nil
}

func GetAllNodes() ([]apiobj.Node, error) {
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

	if res["data"] == nil {
		fmt.Println("empty node list")
		return []apiobj.Node{}, err
	}

	data, ok := res["data"].([]interface{})
	if !ok {
		fmt.Println("expected type []interface{} for field 'data', got something else")
		return nil, fmt.Errorf("type assertion failed for 'data'")
	}

	// 将 interface{} 列表转换为字符串列表
	var nodes []apiobj.Node
	for _, item := range data {
		str, ok := item.(string)
		if !ok {
			fmt.Println("type assertion failed for an item in 'data'")
			return nil, fmt.Errorf("type assertion failed for an item in 'data'")
		}
		var node apiobj.Node
		json.Unmarshal([]byte(str), &node)
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func GetAllReplicaSets() ([]apiobj.ReplicaSet, error) {
	URL := apiconfig.URL_GlobalReplicaSets
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

	if res["data"] == nil {
		fmt.Println("empty replica set list")
		return []apiobj.ReplicaSet{}, err
	}

	data, ok := res["data"].([]interface{})
	if !ok {
		fmt.Println("expected type []interface{} for field 'data', got something else")
		return nil, fmt.Errorf("type assertion failed for 'data'")
	}

	// 将 interface{} 列表转换为字符串列表
	var replicaSets []apiobj.ReplicaSet
	for _, item := range data {
		str, ok := item.(string)
		if !ok {
			fmt.Println("type assertion failed for an item in 'data'")
			return nil, fmt.Errorf("type assertion failed for an item in 'data'")
		}
		var replicaSet apiobj.ReplicaSet
		json.Unmarshal([]byte(str), &replicaSet)
		replicaSets = append(replicaSets, replicaSet)
	}

	return replicaSets, nil
}

func GetAllHpas() ([]apiobj.Hpa, error) {
	URL := apiconfig.URL_GlobalHpas
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

	if res["data"] == nil {
		fmt.Println("empty hpa list")
		return []apiobj.Hpa{}, err
	}

	data, ok := res["data"].([]interface{})
	if !ok {
		fmt.Println("expected type []interface{} for field 'data', got something else")
		return nil, fmt.Errorf("type assertion failed for 'data'")
	}

	// 将 interface{} 列表转换为字符串列表
	var hpas []apiobj.Hpa
	for _, item := range data {
		str, ok := item.(string)
		if !ok {
			fmt.Println("type assertion failed for an item in 'data'")
			return nil, fmt.Errorf("type assertion failed for an item in 'data'")
		}
		var hpa apiobj.Hpa
		json.Unmarshal([]byte(str), &hpa)
		hpas = append(hpas, hpa)
	}

	return hpas, nil
}

func GetAllFunctions () ([]apiobj.Function, error) {
	URL := apiconfig.URL_GlobalFunctions
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

	if res["data"] == nil {
		fmt.Println("empty function list")
		return []apiobj.Function{}, err
	}

	data, ok := res["data"].([]interface{})
	if !ok {
		fmt.Println("expected type []interface{} for field 'data', got something else")
		return nil, fmt.Errorf("type assertion failed for 'data'")
	}

	// 将 interface{} 列表转换为字符串列表
	var functions []apiobj.Function
	for _, item := range data {
		str, ok := item.(string)
		if !ok {
			fmt.Println("type assertion failed for an item in 'data'")
			return nil, fmt.Errorf("type assertion failed for an item in 'data'")
		}
		var function apiobj.Function
		json.Unmarshal([]byte(str), &function)
		functions = append(functions, function)
	}

	return functions, nil
}