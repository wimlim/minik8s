package scheduler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/apirequest"
	"minik8s/pkg/config/apiconfig"
	"minik8s/pkg/message"
	"net/http"
	"strings"

	"github.com/streadway/amqp"
)

func chooseNode() apiobj.Node {
	nodes, err := apirequest.GetAllNodes()
	if err != nil {
		fmt.Println("Error getting nodes:", err)
	}
	return nodes[0]
}

func addPod(msg message.Message) {
	pod := apiobj.Pod{}
	json.Unmarshal([]byte(msg.Content), &pod)

	node := chooseNode()
	pod.Spec.NodeName = node.MetaData.Name

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

	msg_pub := message.Message{
		Type:    "Add",
		URL:     URL,
		Name:    pod.MetaData.Name,
		Content: string(jsonData),
	}
	msgJson, _ := json.Marshal(msg_pub)
	p := message.NewPublisher()
	defer p.Close()

	que := fmt.Sprintf(message.PodQueue+"-%s", node.MetaData.Name)
	p.Publish(que, msgJson)
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
		addPod(msg)
	})
}
