package apiobj

import (
	"time"

	"github.com/docker/docker/api/types"
)

const (
	PodPhase_Pending   = "Pending"
	PodPhase_Running   = "Running"
	PodPhase_Succeeded = "Succeeded"
	PodPhase_Failed    = "Failed"
	PodPhase_Unknown   = "Unknown"
)

type PersistentVolumeClaim struct {
	ClaimName string `yaml:"claimName" json:"claimName"`
}

type VolumeMount struct {
	Name                  string                `yaml:"name" json:"name"`
	MountPath             string                `yaml:"mountPath" json:"mountPath"`
}
type Resource struct {
	CPU    float64 `yaml:"cpu" json:"cpu"`
	Memory float64 `yaml:"memory" json:"memory"`
}
type ContainerPort struct {
	ContainerPort int    `yaml:"containerPort" json:"containerPort"`
	Name          string `yaml:"name" json:"name"`
	Protocol      string `yaml:"protocol" json:"protocol"`
	HostIP        string `yaml:"hostIp" json:"hostIp"`
	HostPort      string `yaml:"hostPort" json:"hostPort"`
}
type Container struct {
	Name         string            `yaml:"name" json:"name"`
	Image        string            `yaml:"image" json:"image"`
	Ports        []ContainerPort   `yaml:"ports" json:"ports"`
	Env          map[string]string `yaml:"env" json:"env"`
	Command      []string          `yaml:"command" json:"command"`
	Args         []string          `yaml:"args" json:"args"`
	Resources    Resource          `yaml:"resources" json:"resources"`
	VolumeMounts []VolumeMount     `yaml:"volumeMounts" json:"volumeMounts"`
	Tty          bool              `yaml:"tty" json:"tty" default:"false"`
}

/*
	支持的 type 值如下：
	取值				行为
	空字符串			（默认）用于向后兼容，这意味着在安装 hostPath 卷之前不会执行任何检查。
	DirectoryOrCreate	如果在给定路径上什么都不存在，那么将根据需要创建空目录，权限设置为 0755，具有与 kubelet 相同的组和属主信息。
	Directory			在给定路径上必须存在的目录。
	FileOrCreate		如果在给定路径上什么都不存在，那么将在那里根据需要创建空文件，权限设置为 0644，具有与 kubelet 相同的组和所有权。
	File				在给定路径上必须存在的文件。
	Socket				在给定路径上必须存在的 UNIX 套接字。
	CharDevice			在给定路径上必须存在的字符设备。
	BlockDevice			在给定路径上必须存在的块设备。
*/

type HostPath struct {
	Path string `yaml:"path" json:"path"`
	Type string `yaml:"type" json:"type"`
}

type Volume struct {
	Name     string   `yaml:"name" json:"name"`
	HostPath HostPath `yaml:"hostPath" json:"hostPath"`
	PersistentVolumeClaim PersistentVolumeClaim `yaml:"persistentVolumeClaim" json:"persistentVolumeClaim"`
}

type PodSpec struct {
	NodeName     string            `yaml:"nodeName" json:"nodeName"`
	Containers   []Container       `yaml:"containers" json:"containers"`
	NodeSelector map[string]string `yaml:"nodeSelector" json:"nodeSelector"`
	Volumes      []Volume          `yaml:"volumes" json:"volumes"`
}

type PodStatus struct {
	Phase          string                 `yaml:"phase" json:"phase"`
	PodIP          string                 `yaml:"podIP" json:"podIP"`
	UpdateTime     time.Time              `yaml:"updateTime" json:"updateTime"`
	CpuUsage       float64                `yaml:"cpuUsage" json:"cpuUsage"`
	MemUsage       float64                `yaml:"memUsage" json:"memUsage"`
	ContainerState []types.ContainerState `yaml:"containerState" json:"containerState"`
}

type Pod struct {
	ApiVersion string    `yaml:"apiVersion" json:"apiVersion"`
	Kind       string    `yaml:"kind" json:"kind"`
	MetaData   MetaData  `yaml:"metadata" json:"metadata"`
	Spec       PodSpec   `yaml:"spec" json:"spec"`
	Status     PodStatus `yaml:"status" json:"status"`
}

func (p *Pod) GetKind() string {
	return p.Kind
}
func (p *Pod) GetName() string {
	return p.MetaData.Name
}
func (p *Pod) GetNamespace() string {
	return p.MetaData.Namespace
}
func (p *Pod) SetNamespace(namespace string) {
	p.MetaData.Namespace = namespace
}
