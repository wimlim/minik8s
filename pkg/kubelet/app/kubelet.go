package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/kubelet/app/runtime"
	"minik8s/pkg/message"
	"os"

	"github.com/streadway/amqp"
)

func msgHandler(d amqp.Delivery) {
	fmt.Println(string(d.Body))
	var msg message.Message
	json.Unmarshal(d.Body, &msg)
	fmt.Println(msg.Name)
	var pod apiobj.Pod
	json.Unmarshal([]byte(msg.Content), &pod)
	runtime.CreatePod(&pod)
	reader := bufio.NewReader(os.Stdin)

	// 读取回车
	fmt.Println("\nPod已创建\n按下回车继续...")
	_, err := reader.ReadBytes('\n')
	if err != nil {
		fmt.Println("读取回车时出错：", err)
		return
	}

	fmt.Println("继续执行下一步操作：删除Pod...")

	runtime.DeletePod(&pod)
	fmt.Println(pod.MetaData.Name)
}

func Run() {
	s := message.NewSubscriber()
	defer s.Close()
	for {
		s.Subscribe(message.PodQueue, msgHandler)
	}
}
