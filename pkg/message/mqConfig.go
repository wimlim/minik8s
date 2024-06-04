package message

import (
	"fmt"
	"minik8s/pkg/config/apiconfig"
)

const (
	RabbitMQUser        = "ling"
	RabbitMQPassword    = "123456"
	RabbitMQDefaultPort = 5672
	DefaultExchange     = "minik8s"
)

func RabbitMQURL() string {
	ip := apiconfig.GetMasterIP()
	url := "amqp://" + RabbitMQUser + ":" + RabbitMQPassword +
		"@" + ip + ":" +
		fmt.Sprint(RabbitMQDefaultPort) + "//"
	return url
}

const (
	ScheduleQueue = "scheduleQueue"
	PodQueue      = "podQueue"
	ServiceQueue  = "serviceQueue"
	DnsQueue      = "dnsQueue"
	WorkflowQueue = "workflowQueue"
)

var Queue2Exchange = map[string]string{
	ScheduleQueue: "minik8s",
	PodQueue:      "minik8s",
	ServiceQueue:  "minik8s",
	DnsQueue:      "minik8s",
	WorkflowQueue: "minik8s",
}
