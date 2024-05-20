package apirequest

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/config/apiconfig"
	"net/http"
	"strings"
)

func GetRequest(namespace string, name string, kind string) (apiobj.ApiObject, error) {

	var apiObject apiobj.ApiObject
	switch kind {
	case "Node":
		apiObject = &apiobj.Node{}
	case "Pod":
		apiObject = &apiobj.Pod{}
	case "Service":
		apiObject = &apiobj.Service{}
	case "ReplicaSet":
		apiObject = &apiobj.ReplicaSet{}
	case "Hpa":
		apiObject = &apiobj.Hpa{}
	case "Dns":
		apiObject = &apiobj.Dns{}
	}

	URL := apiconfig.Kind2URL[kind]
	URL = strings.Replace(URL, ":namespace", namespace, -1)
	URL = strings.Replace(URL, ":name", name, -1)
	URL = apiconfig.GetApiServerUrl() + URL
	response, err := http.Get(URL)
	if err != nil {
		fmt.Printf("get  error")
		return nil, err
	}
	defer response.Body.Close()

	var res map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		fmt.Printf("decode %s error\n", apiObject.GetKind())
		return nil, err
	}
	data := res["data"].(string)

	err = json.Unmarshal([]byte(data), apiObject)
	if err != nil {
		fmt.Printf("unmarshal %s error\n", apiObject.GetKind())
		return nil, err
	}

	return apiObject, nil
}
