package apirequest

import (
	"minik8s/pkg/apiobj"
	"net/http"
	"fmt"
	"encoding/json"
)

func GetRequest(url string, kind string) (apiobj.ApiObject, error) {

	
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
	response, err := http.Get(url)
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