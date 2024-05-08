package minik8sTypes

import "github.com/docker/go-connections/nat"

/*
	Client.ContainerCreate	所需的全部参数
*/

type ContainerConfig struct {
	//	config
	Image        string
	Cmd          []string
	Env          []string
	Tty          bool
	Labels       map[string]string
	Entrypoint   []string
	Volumes      map[string]struct{}
	ExposedPorts map[string]struct{}

	//	host config
	NetworkMode  string
	Binds        []string
	PortBindings nat.PortMap
	IpcMode      string
	PidMode      string
	VolumesFrom  []string
	Links        []string
	Memory       int64
	NanoCPUs     int64

	//	name
	Name string
}

/*
	container	创建需要的const
*/

const (
	Container_Filter_Image = "ancestor"
	Container_Filter_Name  = "name"
	Container_Filter_Id    = "id"
)

const (
	Container_IpcMode_Shareable = "shareable"
)

const (
	Container_Label_PodUid       = "_pod_uid"
	Container_Label_PodName      = "_pod_name"
	Container_Label_IfPause      = "_if_pause"
	Container_Label_PodNamespace = "_pod_namespace"
	Container_Label_IfPause_True = "_true"
)

const (
	Container_Port_Localhost_IP = "127.0.0.1"
	Container_Port_Protocol_TCP = "tcp"
)
