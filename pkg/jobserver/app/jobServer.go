package jobserver

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/config/jobserverconfig"
	"minik8s/pkg/message"

	"github.com/streadway/amqp"
	"golang.org/x/crypto/ssh"
)

func addJob(job apiobj.Job) {

}

func deleteJob(job apiobj.Job) {

}

func Run() {
	sshConfig := &ssh.ClientConfig{
		User: jobserverconfig.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(jobserverconfig.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	sshClient, err := ssh.Dial("tcp", jobserverconfig.HttpUrl, sshConfig)
	if err != nil {
		fmt.Println("Failed to dial: ", err)
	}
	defer sshClient.Close()

	sub := message.NewSubscriber()
	defer sub.Close()
	sub.Subscribe(message.JobQueue, func(d amqp.Delivery) {
		var message message.Message
		err := json.Unmarshal(d.Body, &message)
		if err != nil {
			fmt.Println("Failed to unmarshal message: ", err)
		}
		switch message.Type {
		case "Add":
			fmt.Println("Create job: ")
		case "Delete":
			fmt.Println("Delete job: ")
		}

	})
}
