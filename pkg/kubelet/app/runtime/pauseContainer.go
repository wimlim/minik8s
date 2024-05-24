/*
	对Pause容器的操作
*/

package runtime

import (
	"errors"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/kubelet/app/runtime/container"
	"minik8s/pkg/kubelet/app/runtime/image"
	"minik8s/pkg/minik8sTypes"
	"minik8s/tools/weave"
	"net"
	"strconv"
	"sync"

	"github.com/docker/docker/api/types/filters"
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
		fmt.Println("Error:", err)
		return "", err
	}
	pauseContainerConfig, err := parsePauseContainerConfig(pod)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	for labelKey, labelVal := range pauseContainerConfig.Labels {
		fmt.Println(labelKey + " = " + labelVal)
	}
	pauseId, err := container.CreateContainer(pauseContainerConfig)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	_, err = container.StartContainer(pauseId)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	if pod.Status.PodIP == "" {
		res, err := weave.WeaveAttach(pauseId)
		if err != nil {
			fmt.Println("Error:", err)
			return "", err
		}
		pod.Status.PodIP = res
		fmt.Println("\nPodIP:" + pod.Status.PodIP)
	}
	return pauseId, nil
}

func RemovePauseContainer(pod *apiobj.Pod) (string, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_IfPause+"="+minik8sTypes.Container_Label_IfPause_True)
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_PodUid+"="+pod.MetaData.UID)
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_PodName+"="+pod.MetaData.Name)
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_PodNamespace+"="+pod.MetaData.Namespace)
	// for labelKey, labelValue := range pod.MetaData.Labels {
	// 	filterArgs.Add(minik8sTypes.Container_Filter_Label, labelKey+"="+labelValue)
	// }
	containers, err := container.ListContainerWithFilters(filterArgs)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	if len(containers) == 0 {
		fmt.Println("Pause container has already been removed")
		return "", nil
	} else if len(containers) != 1 {
		return "", errors.New("pause count error")
	}
	_, err = container.RemoveContainer(containers[0].ID)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return containers[0].ID, nil
}

func StartPauseContainer(pod *apiobj.Pod) (string, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_IfPause+"="+minik8sTypes.Container_Label_IfPause_True)
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_PodUid+"="+pod.MetaData.UID)
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_PodName+"="+pod.MetaData.Name)
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_PodNamespace+"="+pod.MetaData.Namespace)
	for labelKey, labelValue := range pod.MetaData.Labels {
		filterArgs.Add(minik8sTypes.Container_Filter_Label, labelKey+"="+labelValue)
	}
	containers, err := container.ListContainerWithFilters(filterArgs)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	if len(containers) != 1 {
		return "", errors.New("pause count error")
	}
	_, err = container.StartContainer(containers[0].ID)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return containers[0].ID, nil
}

func StopPauseContainer(pod *apiobj.Pod) (string, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_IfPause+"="+minik8sTypes.Container_Label_IfPause_True)
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_PodUid+"="+pod.MetaData.UID)
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_PodName+"="+pod.MetaData.Name)
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_PodNamespace+"="+pod.MetaData.Namespace)
	for labelKey, labelValue := range pod.MetaData.Labels {
		filterArgs.Add(minik8sTypes.Container_Filter_Label, labelKey+"="+labelValue)
	}
	containers, err := container.ListContainerWithFilters(filterArgs)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	if len(containers) != 1 {
		return "", errors.New("pause count error")
	}
	_, err = container.StopContainer(containers[0].ID)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return containers[0].ID, nil
}

func RestartPauseContainer(pod *apiobj.Pod) (string, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_IfPause+"="+minik8sTypes.Container_Label_IfPause_True)
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_PodUid+"="+pod.MetaData.UID)
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_PodName+"="+pod.MetaData.Name)
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_PodNamespace+"="+pod.MetaData.Namespace)
	for labelKey, labelValue := range pod.MetaData.Labels {
		filterArgs.Add(minik8sTypes.Container_Filter_Label, labelKey+"="+labelValue)
	}
	containers, err := container.ListContainerWithFilters(filterArgs)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	if len(containers) != 1 {
		return "", errors.New("pause count error")
	}
	_, err = container.RestartContainer(containers[0].ID)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return containers[0].ID, nil
}

func parsePauseContainerConfig(pod *apiobj.Pod) (*minik8sTypes.ContainerConfig, error) {
	// dns
	dns := []string{}
	dns = append(dns, "172.17.0.1")
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
	pauseExposedPorts := map[nat.Port]struct{}{}
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
			bindingPortsKey, err := nat.NewPort(p.Protocol, strconv.Itoa(p.ContainerPort))
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
			pauseExposedPorts[bindingPortsKey] = struct{}{}
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
		Name:         minik8sTypes.Container_Pause_Name_Base + pod.MetaData.UID,
		DNS:          dns,
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
		defer listener.Close()
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
