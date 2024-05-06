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
