/*
	对其他容器(Pause以外的容器)的操作
*/

package runtime

import (
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/kubelet/app/runtime/image"
	"minik8s/pkg/minik8sTypes"
)

func CreateAllCommonContainer(pod *apiobj.Pod, pauseId string) (string, error) {
	for _, commonContainer := range pod.Spec.Containers {
		_, err := createCommonContainer(pod, &commonContainer, pauseId)
		if err != nil {
			RemoveAllPodContainer(pod)
			fmt.Println("Error:", err)
			return "", err
		}
	}
	return "", nil
}

func RemoveAllPodContainer(pod *apiobj.Pod) (string, error) {
	return "", nil
}

func createCommonContainer(pod *apiobj.Pod, commonContainer *apiobj.Container, pauseId string) (string, error) {
	_, err := image.PullImage(commonContainer.Image)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	_, err = parseCommonContainerConfig(pod, commonContainer, pauseId)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return "", nil
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
	podMode := ""
	pauseName := ""

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
		NetworkMode: podMode,
		IpcMode:     podMode,
		PidMode:     podMode,
		Binds:       commonBinds,
		VolumesFrom: []string{pauseName},
		NanoCPUs:    int64(commonContainer.Resources.CPU),
		Memory:      int64(commonContainer.Resources.Memory),
	}
	return &commonContainerConfig, nil
}

func parseVolumeBinds(podVolumes []apiobj.Volume, containerVolumeMounts []apiobj.VolumeMount) ([]string, error) {
	commonBinds := []string{}
	return commonBinds, nil
}
