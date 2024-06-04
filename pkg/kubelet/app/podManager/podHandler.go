package podmanager

import (
	"errors"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/kubelet/app/runtime"
)

type PodHandler struct {
	TaskQueue chan PodTask
}

type PodHandlerInterface interface {
	AddTask(podTask PodTask) error
	Run()
	Stop()
}

func NewPodHandler() *PodHandler {
	return &PodHandler{
		TaskQueue: make(chan PodTask, PodTaskChannelBufferSize),
	}
}

func (p *PodHandler) RunTask(podTask PodTask) error {
	switch podTask.TaskType {
	case Task_AddPod:
		return runtime.CreatePod(podTask.TaskArgs.(*apiobj.Pod))
	case Task_StartPod:
		return runtime.StartPod(podTask.TaskArgs.(*apiobj.Pod))
	case Task_DelPod:
		err := runtime.DeletePod(podTask.TaskArgs.(*apiobj.Pod))
		podTask.OnComplete()
		return err
	case Task_StopPod:
		return runtime.StopPod(podTask.TaskArgs.(*apiobj.Pod))
	case Task_RestartPod:
		return runtime.RestartPod(podTask.TaskArgs.(*apiobj.Pod))
	default:
		return errors.New("unknown task type")
	}
}

func (p *PodHandler) AddTask(podTask PodTask) error {
	if len(p.TaskQueue) == PodTaskChannelBufferSize {
		return errors.New("task queue is full")
	}
	p.TaskQueue <- podTask
	return nil
}

func (p *PodHandler) Run() {
	for task := range p.TaskQueue {
		err := p.RunTask(task)
		if err != nil {
			fmt.Printf("Error in PodHandler: %v\n", err)
		}
	}
}

func (p *PodHandler) Stop() {
	close(p.TaskQueue)
}
