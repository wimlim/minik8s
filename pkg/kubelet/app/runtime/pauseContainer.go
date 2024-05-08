/*
	对Pause容器的操作
*/

package runtime

import (
	"errors"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/kubelet/app/runtime/image"
	"minik8s/pkg/minik8sTypes"
	"net"
	"sync"

	"github.com/docker/go-connections/nat"
)

var lock sync.Mutex

const (
	PauseContainerImageRef = "registry.cn-hangzhou.aliyuncs.com/google_containers/pause:3.6"
)

/*
	CreatePauseContainer
	参数：*apiobj.Pod
	返回：Pause容器Id，error
*/

func CreatePauseContainer(pod *apiobj.Pod) (string, error) {
	_, err := image.PullImage(PauseContainerImageRef)
	if err != nil {
		return "", err
	}
	_, err = parsePauseContainerConfig(pod)
	if err != nil {
		return "", err
	}

	return "", nil
}

func parsePauseContainerConfig(pod *apiobj.Pod) (*minik8sTypes.ContainerConfig, error) {
	//	Labels
	pauseLabels := map[string]string{}
	pauseLabels[minik8sTypes.Container_Label_PodUid] = pod.MetaData.UID
	pauseLabels[minik8sTypes.Container_Label_PodName] = pod.MetaData.Name
	pauseLabels[minik8sTypes.Container_Label_PodNamespace] = pod.MetaData.Namespace
	pauseLabels[minik8sTypes.Container_Label_IfPause] = minik8sTypes.Container_Label_IfPause_True
	for labelKey, labelValue := range pod.MetaData.Labels {
		pauseLabels[labelKey] = labelValue
	}

	// Ports
	pausePortSet := map[string]struct{}{}
	pauseExposedPorts := map[string]struct{}{}
	pauseBindingPorts := nat.PortMap{}
	for _, ctn := range pod.Spec.Containers {
		for _, p := range ctn.Ports {
			if p.HostIP == "" {
				p.HostIP = minik8sTypes.Container_Port_Localhost_IP
			}
			if p.Protocol == "" {
				p.Protocol = minik8sTypes.Container_Port_Protocol_TCP
			}
			if p.HostPort == "" {
				availablePort, err := findAvailablePort(&pausePortSet)
				if err != nil {
					return nil, err
				}
				p.HostPort = availablePort
			}
			pausePortSet[p.HostPort] = struct{}{}
			bindingPortsKey, err := nat.NewPort(p.Protocol, string(p.ContainerPort))
			if err != nil {
				return nil, err
			}
			if _, find := pauseBindingPorts[bindingPortsKey]; find {
				return nil, errors.New("found containerport conflict")
			}
			pauseBindingPorts[bindingPortsKey] = []nat.PortBinding{
				{
					HostIP:   p.HostIP,
					HostPort: p.HostPort,
				},
			}
			pauseExposedPorts[string(bindingPortsKey)] = struct{}{}
		}
	}

	pauseContainerConfig := minik8sTypes.ContainerConfig{
		Image:        PauseContainerImageRef,
		Labels:       pauseLabels,
		ExposedPorts: pauseExposedPorts,
		PortBindings: pauseBindingPorts,
		Volumes:      nil,
		Env:          nil,
		IpcMode:      minik8sTypes.Container_IpcMode_Shareable,
	}
	return &pauseContainerConfig, nil
}

func findAvailablePort(pausePortSet *map[string]struct{}) (string, error) {
	lock.Lock()
	defer lock.Unlock()
	for {
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return "", err
		}
		address := listener.Addr().String()
		_, port, err := net.SplitHostPort(address)
		if err != nil {
			return "", err
		}
		if _, found := (*pausePortSet)[port]; !found {
			return port, nil
		}
	}
}
