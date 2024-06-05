package monitormanager

import "minik8s/pkg/apiobj"

func AddNodeMonitor(node *apiobj.Node) {
	targets := []string{}
	tmpTarget := node.IP + ":10000"
	targets = append(targets, tmpTarget)
	labels := Labels{
		Name:      node.MetaData.Name,
		Namespace: node.MetaData.Namespace,
		UID:       node.MetaData.UID,
	}
	monitorData := MonitorData{
		Targets: targets,
		Labels:  labels,
	}
	AddPodMonitorDataToFile(monitorData)
}

func RemoveNodeMonitor(node *apiobj.Node) {
	labels := Labels{
		Name:      node.MetaData.Name,
		Namespace: node.MetaData.Namespace,
		UID:       node.MetaData.UID,
	}
	RemovePodMonitorDataToFile(labels)
}
