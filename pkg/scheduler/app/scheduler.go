package scheduler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/config/apiconfig"
	"minik8s/pkg/message"
	"net/http"
	"strings"

	"github.com/streadway/amqp"
)

func GetAllNodes() ([]string, error) {
	return getAllNodes()
}

func getAllNodes() ([]string, error) {
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

func chooseNode() string {
	nodes, err := getAllNodes()
	if err != nil {
		fmt.Println("Error getting nodes:", err)
		return ""
	}
	return nodes[0]
}

func addPod(msg message.Message) {
	pod := apiobj.Pod{}
	json.Unmarshal([]byte(msg.Content), &pod)

	node := chooseNode()
	pod.Spec.NodeName = node

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

func deletePod(msg message.Message) {
}

func Run() {
	// subscribe to the schedule queue
	sub := message.NewSubscriber()
	defer sub.Close()
	sub.Subscribe(message.ScheduleQueue, func(d amqp.Delivery) {
		var msg message.Message
		fmt.Println("receive message")
		err := json.Unmarshal(d.Body, &msg)
		if err != nil {
			fmt.Println("unmarshal message error")
			return
		}
		switch msg.Type {
		case "Add":
			fmt.Println("schedule add")
			addPod(msg)
		case "Delete":
			fmt.Println("schedule delete")
			deletePod(msg)
		}
	})
}
