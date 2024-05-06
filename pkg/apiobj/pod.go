package apiobj


type VolumeMount struct {
	Name      string `yaml:"name" json:"name"`
	MountPath string `yaml:"mountPath" json:"mountPath"`
}
type Resource struct {
	CPU    int `yaml:"cpu" json:"cpu"`
	Memory int `yaml:"memory" json:"memory"`
}
type ContainerPort struct {
	ContainerPort int    `yaml:"containerPort" json:"containerPort"`
	Name          string `yaml:"name" json:"name"`
	Protocol      string `yaml:"protocol" json:"protocol"`
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
}
type Volume struct {
	Name string `yaml:"name" json:"name"`
}
type PodSpec struct {
	NodeName     string            `yaml:"nodeName" json:"nodeName"`
	Containers   []Container       `yaml:"containers" json:"containers"`
	NodeSelector map[string]string `yaml:"nodeSelector" json:"nodeSelector"`
	Volumes      []Volume          `yaml:"volumes" json:"volumes"`
}

type PodStatus struct {
	Phase      string `yaml:"phase" json:"phase"`
	PodIP      string `yaml:"podIP" json:"podIP"`
	UpdateTime string `yaml:"updateTime" json:"updateTime"`
}

type Pod struct {
	ApiVersion string      `yaml:"apiVersion" json:"apiVersion"`
	Kind       string      `yaml:"kind" json:"kind"`
	MetaData   MetaData `yaml:"metadata" json:"metadata"`
	Spec       PodSpec     `yaml:"spec" json:"spec"`
	Status     PodStatus   `yaml:"status" json:"status"`
}
