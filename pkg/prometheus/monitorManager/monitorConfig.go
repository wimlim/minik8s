package monitormanager

type Labels struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	UID       string `json:"uuid"`
}

type MonitorData struct {
	Targets []string `json:"targets"`
	Labels  Labels   `json:"labels"`
}

const (
	PodMonitorFilePath = "/srv/prometheus/targets/pods.json"
	NodeMonitorFilePath = "/srv/prometheus/targets/nodes.json"
)