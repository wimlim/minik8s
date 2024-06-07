package apiobj

import "time"

type Node struct {
	ApiVersion string     `yaml:"apiVersion" json:"apiVersion"`
	Kind       string     `yaml:"kind" json:"kind"`
	MetaData   MetaData   `yaml:"metadata" json:"metadata"`
	IP         string     `yaml:"ip" json:"ip"`
	Status     NodeStatus `yaml:"status" json:"status"`
}

type NodeStatus struct {
	Condition  NodeCondition `json:"condition" yaml:"condition"`
	CpuPercent float64       `json:"cpuPercent" yaml:"cpuPercent"`
	MemPercent float64       `json:"memPercent" yaml:"memPercent"`
	PodNum     int           `json:"podNum" yaml:"podNum"`
	UpdateTime time.Time     `json:"updateTime" yaml:"updateTime"`
}

type NodeCondition string

const (
	Ready              NodeCondition = "Ready"
	Unknown            NodeCondition = "Unknown"
	DiskPressure       NodeCondition = "DiskPressure"
	MemoryPressure     NodeCondition = "MemoryPressure"
	PIDPressure        NodeCondition = "PIDPressure"
	NetworkUnavailable NodeCondition = "NetworkUnavailable"
)

func (n *Node) GetKind() string {
	return n.Kind
}
func (n *Node) GetName() string {
	return n.MetaData.Name
}
func (n *Node) GetNamespace() string {
	return n.MetaData.Namespace
}
func (n *Node) SetNamespace(namespace string) {
	n.MetaData.Namespace = namespace
}
