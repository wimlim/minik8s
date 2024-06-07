package status

import (
	"minik8s/pkg/apiobj"
	"minik8s/pkg/kubelet/app/runtime"
	"minik8s/tools/host"
	"time"
)

func GetNodeStatus() (*apiobj.NodeStatus, error) {
	// TODO 完成 node status 实现
	condition := apiobj.NodeCondition(apiobj.Ready)
	cpuPercent, err := host.GetHostCPUPercent()
	if err != nil {
		return nil, err
	}
	memPercent, err := host.GetHostMemoryPercent()
	if err != nil {
		return nil, err
	}
	podNum, err := getNodePodNum()
	if err != nil {
		return nil, err
	}
	nodeStatus := apiobj.NodeStatus{
		Condition:  condition,
		CpuPercent: cpuPercent,
		MemPercent: memPercent,
		PodNum:     podNum,
		UpdateTime: time.Now(),
	}
	return &nodeStatus, nil
}

func getNodePodNum() (int, error) {
	allPodStatus, err := runtime.GetAllPodStatus()
	if err != nil {
		return 0, err
	}
	return len(*allPodStatus), nil
}
