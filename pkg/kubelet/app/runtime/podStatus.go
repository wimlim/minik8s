package runtime

import (
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/kubelet/app/runtime/container"
	"minik8s/pkg/minik8sTypes"
	"minik8s/tools/weave"
	"time"

	"github.com/docker/docker/api/types"
)

func GetAllPodStatus() (*map[minik8sTypes.PodIdentifier]*apiobj.PodStatus, error) {
	containers, err := container.ListAllContainer()
	if err != nil {
		return nil, err
	}
	mapPodId2ContainerIds := map[minik8sTypes.PodIdentifier][]string{}
	mapPodId2PodStatus := map[minik8sTypes.PodIdentifier]*apiobj.PodStatus{}
	for _, ctn := range containers {
		if ctn.Labels[minik8sTypes.Container_Label_PodUid] == "" ||
			ctn.Labels[minik8sTypes.Container_Label_PodName] == "" ||
			ctn.Labels[minik8sTypes.Container_Label_PodNamespace] == "" {
			continue
		}
		podIdentifier := minik8sTypes.PodIdentifier{
			PodId:        ctn.Labels[minik8sTypes.Container_Label_PodUid],
			PodName:      ctn.Labels[minik8sTypes.Container_Label_PodName],
			PodNamespace: ctn.Labels[minik8sTypes.Container_Label_PodNamespace],
		}
		mapPodId2ContainerIds[podIdentifier] = append(mapPodId2ContainerIds[podIdentifier], ctn.ID)
	}

	for podIdentifier, containerIds := range mapPodId2ContainerIds {
		mapPodId2PodStatus[podIdentifier], err = getPodStatus(podIdentifier.PodId, &containerIds)
		if err != nil {
			return nil, err
		}
	}
	return &mapPodId2PodStatus, nil
}

func getPodStatus(podId string, containerIds *[]string) (*apiobj.PodStatus, error) {
	containerStates := []types.ContainerState{}
	containerIPs := []string{}
	podCpuUsage := float64(0.0)
	podMemUsage := float64(0.0)
	for _, ctnId := range *containerIds {
		inspectRes, err := container.InspectContainer(ctnId)
		if err != nil {
			return nil, err
		}
		if inspectRes == nil {
			containerStates = append(containerStates, types.ContainerState{})
		} else {
			containerStates = append(containerStates, *inspectRes.State)
		}

		containerIP, err := weave.WeaveFindIpByContainerID(ctnId)
		if err != nil {
			return nil, err
		}
		containerIPs = append(containerIPs, containerIP)

		containerCpuUsage, containerMemUsage, err := container.CalcContainerCPUAndMemoryUsage(ctnId)
		if err != nil {
			return nil, err
		}
		podCpuUsage += containerCpuUsage
		podMemUsage += containerMemUsage
	}

	podIP := ""
	if len(containerIPs) != 0 {
		for _, containerIP := range containerIPs {
			if containerIP != containerIPs[0] {
				fmt.Println("getPodStatus: pod " + podId + " IP is no the same")
				break
			}
		}
		podIP = containerIPs[0]
	}
	podPhase, err := parsePodPhaseByContainerStates(&containerStates)
	if err != nil {
		fmt.Println("parsePodPhaseByContainerStates has something wrong")
	}
	podStatus := apiobj.PodStatus{
		Phase:          podPhase,
		PodIP:          podIP,
		UpdateTime:     time.Now(),
		CpuUsage:       podCpuUsage,
		MemUsage:       podMemUsage,
		ContainerState: containerStates,
	}
	return &podStatus, nil
}

func parsePodPhaseByContainerStates(containerStates *[]types.ContainerState) (string, error) {
	if len(*containerStates) == 0 {
		return apiobj.PodPhase_Pending, nil
	}
	isAllCreated := true
	hasRunning := false
	isAllDead := true
	hasFailedContainer := false

	for _, containerState := range *containerStates {
		isAllCreated = isAllCreated && (containerState.Status == "created")
		hasRunning = hasRunning || containerState.Running
		isAllDead = isAllDead && containerState.Dead
		hasFailedContainer = hasFailedContainer || (containerState.ExitCode != 0)
	}

	if isAllCreated {
		return apiobj.PodPhase_Pending, nil
	}

	if hasFailedContainer {
		return apiobj.PodPhase_Failed, nil
	}

	if isAllDead && !hasFailedContainer {
		return apiobj.PodPhase_Succeeded, nil
	}

	if hasRunning {
		return apiobj.PodPhase_Running, nil
	}

	return apiobj.PodPhase_Unknown, nil
}
