package podmanager

import "fmt"

type PodHandler struct {
	TaskQueue chan PodTask
}

type PodHandlerInterface interface {
	AddTask(podTask PodTask) error
	RunTask(podTask PodTask) error
	Run()
	Stop()
}

func NewPodHandler() *PodHandler {
	return &PodHandler{
		TaskQueue: make(chan PodTask, PodTaskChannelBufferSize),
	}
}

func (p *PodHandler) RunTask(podTask PodTask) error {
	// TODO Finish PodHandler
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
