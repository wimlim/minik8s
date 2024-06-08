package monitormanager

import "minik8s/pkg/apiobj"

func AddPodMonitor(pod *apiobj.Pod) {
	if len(pod.Monitor.MetricsPorts) == 0 {
		return
	}
	targets := []string{}
	metricsPorts := pod.Monitor.MetricsPorts
	for _, metricsPort := range metricsPorts {
		tmpTarget := pod.Status.PodIP + ":" + metricsPort
		targets = append(targets, tmpTarget)
	}
	labels := Labels{
		Name:      pod.MetaData.Name,
		Namespace: pod.MetaData.Namespace,
		UID:       pod.MetaData.UID,
	}
	monitorData := MonitorData{
		Targets: targets,
		Labels:  labels,
	}
	AddPodMonitorDataToFile(monitorData)
}

func RemovePodMonitor(pod *apiobj.Pod) {
	if len(pod.Monitor.MetricsPorts) == 0 {
		return
	}
	labels := Labels{
		Name:      pod.MetaData.Name,
		Namespace: pod.MetaData.Namespace,
		UID:       pod.MetaData.UID,
	}
	RemovePodMonitorDataToFile(labels)
}
