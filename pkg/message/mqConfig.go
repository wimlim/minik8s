package message
import (
	"minik8s/pkg/config/apiconfig"
	"fmt"
)
const (
	RabbitMQUser = "ling"
	RabbitMQPassword = "123456"
	RabbitMQDefaultPort = 5672
)
func RabbitMQURL() string {
	ip := apiconfig.GetMasterIP()
	url := "amqp://" + RabbitMQUser + ":" + RabbitMQPassword + 
	"@" + ip + ":" + 
	fmt.Sprint(RabbitMQDefaultPort) + "//"
	return url
}
const(
	ScheduleQueue = "scheduleQueue"
	PodQueue = "podQueue"
)
var Queue2Exchange = map[string]string{
	ScheduleQueue: "minik8s",
	PodQueue: "minik8s",
}