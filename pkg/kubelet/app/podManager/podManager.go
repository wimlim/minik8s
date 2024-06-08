/*
	对pod实现并发操作
*/

package podmanager

import (
	"errors"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/kubelet/app/runtime"
	"time"
)

type PodManager struct {
	PodHandlersMap map[string]*PodHandler
}

type PodManagerInterface interface {
	AddPod(pod *apiobj.Pod) error
	StartPod(pod *apiobj.Pod) error
	DeletePod(pod *apiobj.Pod) error
	StopPod(pod *apiobj.Pod) error
	RestartPod(pod *apiobj.Pod) error
}

func NewPodManager() *PodManager {
	podHandlerMap := map[string]*PodHandler{}
	allPodStatus, err := runtime.GetAllPodStatus()
	if err != nil {
		fmt.Println(err.Error())
	}
	for podIdentifier := range *allPodStatus {
		podHandlerMap[podIdentifier.PodId] = NewPodHandler()
		go podHandlerMap[podIdentifier.PodId].Run()
	}
	return &PodManager{
		PodHandlersMap: podHandlerMap,
	}
}

func (p *PodManager) AddPod(pod *apiobj.Pod) error {
	podId := pod.MetaData.UID
	if _, ok := p.PodHandlersMap[podId]; ok {
		return errors.New("pod already exists")
	}
	podHandler := NewPodHandler()
	p.PodHandlersMap[podId] = podHandler
	go podHandler.Run()

	task := PodTask{
		TaskType: Task_AddPod,
		TaskArgs: pod,
	}
	err := podHandler.AddTask(task)
	time.Sleep(1 * time.Second)
	if err != nil {
		return err
	}
	return nil
}

func (p *PodManager) StartPod(pod *apiobj.Pod) error {
	podId := pod.MetaData.UID
	if _, ok := p.PodHandlersMap[podId]; !ok {
		return errors.New("pod not exists")
	}
	task := PodTask{
		TaskType: Task_StartPod,
		TaskArgs: pod,
	}
	err := p.PodHandlersMap[podId].AddTask(task)
	time.Sleep(1 * time.Second)
	if err != nil {
		return err
	}
	return nil
}

func (p *PodManager) StopPod(pod *apiobj.Pod) error {
	podId := pod.MetaData.UID
	if _, ok := p.PodHandlersMap[podId]; !ok {
		return errors.New("pod not exists")
	}
	task := PodTask{
		TaskType: Task_StopPod,
		TaskArgs: pod,
	}
	err := p.PodHandlersMap[podId].AddTask(task)
	time.Sleep(1 * time.Second)
	if err != nil {
		return err
	}
	return nil
}

func (p *PodManager) RestartPod(pod *apiobj.Pod) error {
	podId := pod.MetaData.UID
	if _, ok := p.PodHandlersMap[podId]; !ok {
		return errors.New("pod not exists")
	}
	task := PodTask{
		TaskType: Task_RestartPod,
		TaskArgs: pod,
	}
	err := p.PodHandlersMap[podId].AddTask(task)
	time.Sleep(1 * time.Second)
	if err != nil {
		return err
	}
	return nil
}

func (p *PodManager) DeletePod(pod *apiobj.Pod) error {
	podId := pod.MetaData.UID
	if _, ok := p.PodHandlersMap[podId]; !ok {
		return errors.New("pod not exists")
	}
	done := make(chan struct{})
	task := PodTask{
		TaskType: Task_DelPod,
		TaskArgs: pod,
		OnComplete: func() {
			close(done)
		},
	}
	err := p.PodHandlersMap[podId].AddTask(task)
	if err != nil {
		return err
	}
	<-done
	p.PodHandlersMap[podId].Stop()
	delete(p.PodHandlersMap, podId)
	return nil
}
