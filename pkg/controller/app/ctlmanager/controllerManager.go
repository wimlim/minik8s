package ctlmanager

import (
	"minik8s/pkg/controller/app/controllers"
)	
type ControllerManager struct {
	rc *controllers.ReplicaController
	hc *controllers.HpaController
}

func NewControllerManager() *ControllerManager {
	return &ControllerManager{
		rc: controllers.NewReplicaController(),
		hc: controllers.NewHpaController(),
	}
}

func (cm *ControllerManager) Run(stop <-chan struct{}) {
	go cm.rc.Run()
	go cm.hc.Run()

	_, ok := <-stop
	if ok {
		return
	}
}