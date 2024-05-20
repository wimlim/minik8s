/*
	对其他容器(Pause以外的容器)的操作
*/

package runtime

import (
	"errors"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/kubelet/app/runtime/container"
	"minik8s/pkg/kubelet/app/runtime/image"
	"minik8s/pkg/minik8sTypes"

	"github.com/docker/docker/api/types/filters"
	"github.com/google/uuid"
)

func CreateAllCommonContainer(pod *apiobj.Pod, pauseId string) (string, error) {
	for _, commonContainer := range pod.Spec.Containers {
		_, err := createCommonContainer(pod, &commonContainer, pauseId)
		if err != nil {
			RemoveAllCommonContainer(pod)
			fmt.Println("Error:", err)
			return "", err
		}
	}
	return pod.MetaData.UID, nil
}

func RemoveAllCommonContainer(pod *apiobj.Pod) (string, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_IfPause+"="+minik8sTypes.Container_Label_IfPause_False)
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_PodUid+"="+pod.MetaData.UID)
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_PodName+"="+pod.MetaData.Name)
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_PodNamespace+"="+pod.MetaData.Namespace)
	// for labelKey, labelValue := range pod.MetaData.Labels {
	// 	filterArgs.Add(minik8sTypes.Container_Filter_Label, labelKey+"="+labelValue)
	// }
	commonContainers, err := container.ListContainerWithFilters(filterArgs)
	if err != nil {
		fmt.Println("Error in RemoveAllCommonContainer: container.ListContainerWithFilters")
		return "", err
	}
	errMsg := ""
	for _, commonContainer := range commonContainers {
		_, err = container.RemoveContainer(commonContainer.ID)
		if err != nil {
			errMsg += "RemoveContainer error in container_" + commonContainer.ID + ":\n" + err.Error() + "\n"
		}
	}
	if errMsg != "" {
		return "", errors.New(errMsg)
	}
	return pod.MetaData.UID, nil
}

func StartAllCommonContainer(pod *apiobj.Pod) (string, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_IfPause+"="+minik8sTypes.Container_Label_IfPause_False)
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_PodUid+"="+pod.MetaData.UID)
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_PodName+"="+pod.MetaData.Name)
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_PodNamespace+"="+pod.MetaData.Namespace)
	for labelKey, labelValue := range pod.MetaData.Labels {
		filterArgs.Add(minik8sTypes.Container_Filter_Label, labelKey+"="+labelValue)
	}
	commonContainers, err := container.ListContainerWithFilters(filterArgs)
	if err != nil {
		fmt.Println("Error in StartAllCommonContainer: container.ListContainerWithFilters")
		return "", err
	}
	errMsg := ""
	for _, commonContainer := range commonContainers {
		_, err = container.StartContainer(commonContainer.ID)
		if err != nil {
			errMsg += "StartContainer error in container_" + commonContainer.ID + ":\n" + err.Error() + "\n"
		}
	}
	if errMsg != "" {
		return "", errors.New(errMsg)
	}
	return pod.MetaData.UID, nil
}

func StopAllCommonContainer(pod *apiobj.Pod) (string, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_IfPause+"="+minik8sTypes.Container_Label_IfPause_False)
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_PodUid+"="+pod.MetaData.UID)
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_PodName+"="+pod.MetaData.Name)
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_PodNamespace+"="+pod.MetaData.Namespace)
	for labelKey, labelValue := range pod.MetaData.Labels {
		filterArgs.Add(minik8sTypes.Container_Filter_Label, labelKey+"="+labelValue)
	}
	commonContainers, err := container.ListContainerWithFilters(filterArgs)
	if err != nil {
		fmt.Println("Error in StopAllCommonContainer: container.ListContainerWithFilters")
		return "", err
	}
	errMsg := ""
	for _, commonContainer := range commonContainers {
		_, err = container.StopContainer(commonContainer.ID)
		if err != nil {
			errMsg += "StopContainer error in container_" + commonContainer.ID + ":\n" + err.Error() + "\n"
		}
	}
	if errMsg != "" {
		return "", errors.New(errMsg)
	}
	return pod.MetaData.UID, nil
}

func RestartAllCommonContainer(pod *apiobj.Pod) (string, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_IfPause+"="+minik8sTypes.Container_Label_IfPause_False)
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_PodUid+"="+pod.MetaData.UID)
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_PodName+"="+pod.MetaData.Name)
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_PodNamespace+"="+pod.MetaData.Namespace)
	for labelKey, labelValue := range pod.MetaData.Labels {
		filterArgs.Add(minik8sTypes.Container_Filter_Label, labelKey+"="+labelValue)
	}
	commonContainers, err := container.ListContainerWithFilters(filterArgs)
	if err != nil {
		fmt.Println("Error in RestartAllCommonContainer: container.ListContainerWithFilters")
		return "", err
	}
	errMsg := ""
	for _, commonContainer := range commonContainers {
		_, err = container.RestartContainer(commonContainer.ID)
		if err != nil {
			errMsg += "RestartContainer error in container_" + commonContainer.ID + ":\n" + err.Error() + "\n"
		}
	}
	if errMsg != "" {
		return "", errors.New(errMsg)
	}
	return pod.MetaData.UID, nil
}

func createCommonContainer(pod *apiobj.Pod, commonContainer *apiobj.Container, pauseId string) (string, error) {
	_, err := image.PullImage(commonContainer.Image)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	commonContainerConfig, err := parseCommonContainerConfig(pod, commonContainer, pauseId)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	commonContainerId, err := container.CreateContainer(commonContainerConfig)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	_, err = container.StartContainer(commonContainerId)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return commonContainerId, nil
}

func parseCommonContainerConfig(pod *apiobj.Pod, commonContainer *apiobj.Container, pauseId string) (*minik8sTypes.ContainerConfig, error) {
	//	Labels
	commonLabels := map[string]string{}
	commonLabels[minik8sTypes.Container_Label_PodUid] = pod.MetaData.UID
	commonLabels[minik8sTypes.Container_Label_PodName] = pod.MetaData.Name
	commonLabels[minik8sTypes.Container_Label_PodNamespace] = pod.MetaData.Namespace
	commonLabels[minik8sTypes.Container_Label_IfPause] = minik8sTypes.Container_Label_IfPause_False
	for labelKey, labelValue := range pod.MetaData.Labels {
		commonLabels[labelKey] = labelValue
	}
	//	Env
	commonEnv := []string{}
	for envKey, envValue := range commonContainer.Env {
		commonEnv = append(commonEnv, envKey+"="+envValue)
	}
	//	Binds
	commonBinds, err := parseVolumeBinds(pod.Spec.Volumes, commonContainer.VolumeMounts)
	if err != nil {
		return nil, err
	}

	pauseRef := minik8sTypes.Pause_Container_Ref + pauseId
	pauseName := minik8sTypes.Pause_Container_Namebase + pod.MetaData.UID

	if commonContainer.Name == "" {
		commonContainer.Name = minik8sTypes.Common_Container_Namebase + NewUUID()
	}
	commonContainerConfig := minik8sTypes.ContainerConfig{
		// config
		Image:      commonContainer.Image,
		Cmd:        commonContainer.Args,
		Env:        commonEnv,
		Tty:        commonContainer.Tty,
		Labels:     commonLabels,
		Entrypoint: commonContainer.Command,
		Volumes:    nil,
		// host config
		NetworkMode: pauseRef,
		IpcMode:     pauseRef,
		PidMode:     pauseRef,
		Binds:       commonBinds,
		VolumesFrom: []string{pauseName},
		NanoCPUs:    int64(commonContainer.Resources.CPU),
		Memory:      int64(commonContainer.Resources.Memory),
		//	name
		Name: commonContainer.Name,
	}
	return &commonContainerConfig, nil
}

func parseVolumeBinds(podVolumes []apiobj.Volume, containerVolumeMounts []apiobj.VolumeMount) ([]string, error) {
	commonBinds := []string{}
	name_podVolume_map := map[string]*apiobj.Volume{}
	for _, podVolume := range podVolumes {
		if podVolume.HostPath.Path != "" {
			name_podVolume_map[podVolume.Name] = &podVolume
		}
	}
	for _, containerVolumeMount := range containerVolumeMounts {
		_, find := name_podVolume_map[containerVolumeMount.Name]
		if !find {
			return []string{}, errors.New("cannot find container volume")
		}
		pV := name_podVolume_map[containerVolumeMount.Name]
		if pV.HostPath.Path == "" {
			return []string{}, errors.New("the hostpath is empty or undefined")
		}
		commonBinds = append(commonBinds, pV.HostPath.Path+":"+containerVolumeMount.MountPath)
	}
	return commonBinds, nil
}

func NewUUID() string {
	uuid := uuid.New()
	return fmt.Sprintf("NewUUID:%s\n", uuid)
}
