/*
	对Pause容器的操作
*/

package runtime

import (
	"minik8s/pkg/apiobj"
	"minik8s/pkg/kubelet/app/runtime/image"
	"minik8s/pkg/minik8sTypes"

	"github.com/docker/go-connections/nat"
)

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
	pauseContainerConfig, err := parsePauseContainerConfig(pod)
	if err != nil {
		return "", err
	}
	return "", nil
}

func parsePauseContainerConfig(pod *apiobj.Pod) (minik8sTypes.ContainerConfig, error) {
	pauseLabels := map[string]string{}
	pauseExposedPorts := map[string]struct{}{}
	pauseBindingPorts := nat.PortMap{}
	pauseContainerConfig := minik8sTypes.ContainerConfig{
		Image:        PauseContainerImageRef,
		Labels:       pauseLabels,
		ExposedPorts: pauseExposedPorts,
		PortBindings: pauseBindingPorts,
		Volumes:      nil,
		Env:          nil,
		IpcMode:      minik8sTypes.Container_IpcMode_Shareable,
	}
	return pauseContainerConfig, nil
}
