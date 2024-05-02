package message
import (
	"minik8s/pkg/config/apiconfig"
	"fmt"
)
func RabbitMQDefaultURL() string {

	url := "amqp://ling:123456@" + apiconfig.ServerLocalIP + ":" + 
	fmt.Sprint(5672) + "//"
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