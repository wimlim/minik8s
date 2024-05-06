package scheduler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/config/apiconfig"
	"minik8s/pkg/message"
	"net/http"
	"strings"

	"github.com/streadway/amqp"
)

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

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	var nodes []string
	err = json.Unmarshal(body, &nodes)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return nil, err
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
	response, err := http.Post(HttpUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("post error")
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
