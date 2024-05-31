package status

import (
	"errors"
	"fmt"
	"minik8s/pkg/apiobj"
	apiserverutil "minik8s/pkg/kubelet/app/apiserverUtil"
	"minik8s/pkg/kubelet/app/cache"
	"minik8s/pkg/kubelet/app/runtime"
	"minik8s/tools/runner"
	"time"
)

func pushNodeStatus() error {
	// TODO 完成 Node Status 实现
	// nodeStatus, err := GetNodeStatus()
	// if err != nil {
	// 	return err
	// }
	return nil
}

func pushAllPodStatus() error {
	allPodStatus, err := runtime.GetAllPodStatus()
	for podIdentifier, podStatus := range *allPodStatus {
		apiserverutil.PodStatusUpdate(podIdentifier, podStatus)
	}
	return err
}

func pullAllPodStatus(c *cache.PodCache) ([]apiobj.Pod, error) {
	allRemotePods, err := apiserverutil.GetAllRemotePods()
	if err != nil {
		return nil, err
	}
	remotePodsMap, err := convertRemotePodsMap(allRemotePods)
	if err != nil {
		return nil, err
	}
	err = updatePodStatusInCache(remotePodsMap, c)
	if err != nil {
		return nil, err
	}
	return allRemotePods, nil
}

func convertRemotePodsMap(allRemotePods []apiobj.Pod) (map[string]*apiobj.Pod, error) {
	remotePodsMap := map[string]*apiobj.Pod{}
	for _, pod := range allRemotePods {
		podId := pod.MetaData.UID
		remotePodsMap[podId] = &pod
	}
	return remotePodsMap, nil
}

func updatePodStatusInCache(remotePodsMap map[string]*apiobj.Pod, c *cache.PodCache) error {
	localPodsMap, err := c.GetAllPodFromCache()
	if err != nil {
		return err
	}

	errInfo := ""

	for localPodId, localPod := range localPodsMap {
		if _, ok := remotePodsMap[localPodId]; !ok {
			err := c.DeletePodFromCache(localPodId)
			if err != nil {
				errInfo += err.Error() + "\n"
			}
		} else {
			remotePod := remotePodsMap[localPodId]
			if localPod.Status.UpdateTime.Before(remotePod.Status.UpdateTime) {
				err := c.UpdatePodFromCache(remotePod)
				if err != nil {
					errInfo += err.Error() + "\n"
				}
			}
		}
	}

	for remotePodId, remotePod := range remotePodsMap {
		if _, ok := localPodsMap[remotePodId]; !ok {
			err := c.UpdatePodFromCache(remotePod)
			if err != nil {
				errInfo += err.Error() + "\n"
			}
		}
	}

	if errInfo != "" {
		return errors.New(errInfo)
	}

	return nil
}

func Run(c *cache.PodCache) {
	r := runner.NewRunner()
	runPushNodeStatus := func() {
		err := pushNodeStatus()
		if err != nil {
			fmt.Println("error: " + err.Error())
		}
	}
	runPushAllPodStatus := func() {
		err := pushAllPodStatus()
		if err != nil {
			fmt.Println("error: " + err.Error())
		}
	}
	runPullAllPodStatus := func() {
		_, err := pullAllPodStatus(c)
		if err != nil {
			fmt.Println("error: " + err.Error())
		}
	}
	go r.RunLoop(0*time.Second, 5*time.Second, runPushNodeStatus)
	go r.RunLoop(0*time.Second, 5*time.Second, runPushAllPodStatus)
	go r.RunLoop(0*time.Second, 1000*time.Second, runPullAllPodStatus)
}
