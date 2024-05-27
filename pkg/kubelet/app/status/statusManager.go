package status

import (
	"fmt"
	"minik8s/pkg/kubelet/app/runtime"
	"minik8s/tools/runner"
	"time"
)

func pushAllPodStatus() error {
	allPodStatus, err := runtime.GetAllPodStatus()
	if err != nil {
		return err
	}
	for _, _ = range *allPodStatus {

	}
	return nil
}

func run() {
	runPushAllPodStatus := func() {
		err := pushAllPodStatus()
		if err != nil {
			fmt.Println("error: " + err.Error())
		}
	}
	go runner.NewRunner().RunLoop(0*time.Second, 5*time.Second, runPushAllPodStatus)
}
