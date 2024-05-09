package minik8sTypes

import "github.com/docker/go-connections/nat"

/*
	Client.ContainerCreate	所需的全部参数
*/

type ContainerConfig struct {
	// config
	Image        string                `json:"image"`
	Cmd          []string              `json:"cmd"`
	Env          []string              `json:"env"`
	Tty          bool                  `json:"tty"`
	Labels       map[string]string     `json:"labels"`
	Entrypoint   []string              `json:"entrypoint"`
	Volumes      map[string]struct{}   `json:"volumes"`
	ExposedPorts map[nat.Port]struct{} `json:"exposedPorts"`

	// host config
	NetworkMode  string      `json:"networkMode"`
	Binds        []string    `json:"binds"`
	PortBindings nat.PortMap `json:"portBindings"`
	IpcMode      string      `json:"ipcMode"`
	PidMode      string      `json:"pidMode"`
	VolumesFrom  []string    `json:"volumesFrom"`
	Links        []string    `json:"links"`
	Memory       int64       `json:"memory"`
	NanoCPUs     int64       `json:"nanoCPUs"`

	// name
	Name string `json:"name"`
}

/*
	container	创建需要的const
*/

const (
	Container_Filter_Image = "ancestor"
	Container_Filter_Name  = "name"
	Container_Filter_Id    = "id"
	Container_Filter_Label = "label"
)

const (
	Container_IpcMode_Shareable = "shareable"
)

const (
	Container_Label_PodUid        = "_pod_uid"
	Container_Label_PodName       = "_pod_name"
	Container_Label_IfPause       = "_if_pause"
	Container_Label_PodNamespace  = "_pod_namespace"
	Container_Label_IfPause_True  = "_true"
	Container_Label_IfPause_False = "_false"
)

const (
	Container_Port_Localhost_IP = "127.0.0.1"
	Container_Port_Protocol_TCP = "tcp"
)

const (
	Container_Pause_Name_Base = "pause-"
)
