package runtime

import (
	"fmt"
	"io/ioutil"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/kubelet/app/runtime/container"
	"minik8s/pkg/minik8sTypes"
	"testing"

	"github.com/docker/docker/api/types/filters"
	"gopkg.in/yaml.v3"
)

func getPodFromYaml(filePath string) (*apiobj.Pod, error) {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("读取 YAML 文件失败：%v", err)
		return nil, err
	}
	var pod apiobj.Pod
	err = yaml.Unmarshal(yamlFile, &pod)
	if err != nil {
		fmt.Printf("解析 YAML 文件失败：%v", err)
		return nil, err
	}
	return &pod, nil
}

func TestFindAvailablePort(t *testing.T) {
	pausePortSet := map[string]struct{}{}
	for i := 1; i <= 100; i++ {
		port, err := findAvailablePort(&pausePortSet)
		if err != nil {
			fmt.Printf("Error finding available port: %v\n", err)
			return
		}
		pausePortSet[port] = struct{}{}
		fmt.Printf("Test %d: Available Non-repeating port: %s\n", i, port)
	}
	fmt.Printf("\npausePortSet:\n")
	for port, _ := range pausePortSet {
		fmt.Printf("%s\n", port)
	}
}

func TestGetPodFromYaml(t *testing.T) {
	pod, _ := getPodFromYaml("pod1.yaml")
	fmt.Printf("Pod 名称: %s\n", pod.MetaData.Name)
	fmt.Println("容器信息:")
	for _, container := range pod.Spec.Containers {
		fmt.Printf("- 名称: %s, 镜像: %s\n", container.Name, container.Image)
	}
}

func TestRemoveAllCommonContainer(t *testing.T) {
	pod, _ := getPodFromYaml("pod1.yaml")
	filterArgs := filters.NewArgs()
	filterArgs.Add(minik8sTypes.Container_Filter_Label, minik8sTypes.Container_Label_PodName+"="+pod.MetaData.Name)
	for labelKey, labelValue := range pod.MetaData.Labels {
		filterArgs.Add(minik8sTypes.Container_Filter_Label, labelKey+"="+labelValue)
	}
	commonContainers, _ := container.ListContainerWithFilters(filterArgs)
	for _, ctn := range commonContainers {
		fmt.Println(ctn.ID + "\t" + ctn.Image + "\t" + ctn.Names[0])
	}
	errMsg := ""
	for _, commonContainer := range commonContainers {
		_, err := container.RemoveContainer(commonContainer.ID)
		if err != nil {
			errMsg += "RemoveContainer error in container_" + commonContainer.ID + ":\n" + err.Error() + "\n"
		}
	}
	if errMsg != "" {
		return
	}
	return
}
