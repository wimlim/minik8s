package podmanager

type PodTask struct {
	TaskType string
}

const (
	Task_None       = "None"
	Task_AddPod     = "AddPod"
	Task_StartPod   = "StartPod"
	Task_DelPod     = "DelPod"
	Task_StopPod    = "StopPod"
	Task_RestartPod = "RestartPod"
)

const (
	PodTaskChannelBufferSize = 20
)
