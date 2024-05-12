package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/config/apiconfig"
	"minik8s/pkg/kubelet/app/runtime"
	"minik8s/pkg/message"
	"net/http"
	"strings"

	"github.com/streadway/amqp"
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

func msgHandler(d amqp.Delivery) {
	fmt.Println(string(d.Body))
	var msg message.Message
	json.Unmarshal(d.Body, &msg)
	fmt.Println(msg.Name)
	var pod apiobj.Pod
	json.Unmarshal([]byte(msg.Content), &pod)
	if msg.Type == "Delete" {
		runtime.DeletePod(&pod)
		fmt.Println(pod.MetaData.Name)
	} else if msg.Type == "Add" {
		runtime.CreatePod(&pod)
		fmt.Println(pod.MetaData.Name)
		PodUpdate(&pod)
	}
}

func Run() {
	s := message.NewSubscriber()
	defer s.Close()
	for {
		s.Subscribe(message.PodQueue, msgHandler)
	}
}
