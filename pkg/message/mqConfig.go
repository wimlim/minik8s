package message

import (
	"fmt"
	"minik8s/pkg/config/apiconfig"
)

const (
	RabbitMQUser        = "ling"
	RabbitMQPassword    = "123456"
	RabbitMQDefaultPort = 5672
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
	JobQueue      = "jobQueue"
)

var Queue2Exchange = map[string]string{
	ScheduleQueue: "minik8s",
	PodQueue:      "minik8s",
	ServiceQueue:  "minik8s",
	DnsQueue:      "minik8s",
	JobQueue:      "minik8s",
}
