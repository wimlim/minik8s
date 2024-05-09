/*
实现Pod抽象
*/
package runtime

import "minik8s/pkg/apiobj"

func CreatePod(pod *apiobj.Pod) error {
	_, err := CreatePauseContainer(pod)
	if err != nil {
		return err
	}
	_, err = CreateAllCommonContainer(pod)
	if err != nil {
		return err
	}
	return nil
}
