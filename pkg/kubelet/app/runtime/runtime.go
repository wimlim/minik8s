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

func DeletePod(pod *apiobj.Pod) error {
	_, err := RemoveAllCommonContainer(pod)
	if err != nil {
		return err
	}
	_, err = RemovePauseContainer(pod)
	if err != nil {
		return err
	}
	return nil
}

func StopPod(pod *apiobj.Pod) error {
	_, err := StopAllCommonContainer(pod)
	if err != nil {
		return err
	}
	_, err = StopPauseContainer(pod)
	if err != nil {
		return err
	}
	return nil
}

func RestartPod(pod *apiobj.Pod) error {
	_, err := RestartPauseContainer(pod)
	if err != nil {
		return err
	}
	_, err = RestartAllCommonContainer(pod)
	if err != nil {
		return err
	}
	return nil
}
