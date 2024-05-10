/*
实现Pod抽象
*/
package runtime

import "minik8s/pkg/apiobj"

func CreatePod(pod *apiobj.Pod) error {
	pauseId, err := CreatePauseContainer(pod)
	if err != nil {
		return err
	}
	_, err = CreateAllCommonContainer(pod, pauseId)
	if err != nil {
		return err
	}
	return nil
}

func StartPod(pod *apiobj.Pod) error {
	_, err := StartPauseContainer(pod)
	if err != nil {
		return err
	}
	_, err = StartAllCommonContainer(pod)
	if err != nil {
		return err
	}
	return nil
}
