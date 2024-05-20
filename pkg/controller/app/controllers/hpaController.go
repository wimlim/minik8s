package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/apirequest"
	"minik8s/pkg/config/apiconfig"
	"minik8s/tools/runner"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type HpaController struct {
}

func NewHpaController() *HpaController {
	return &HpaController{}
}

func (hc *HpaController) Run() {
	rr := runner.NewRunner()
	rr.RunLoop(5*time.Second, 5*time.Second, hc.update_hpa_pod)
}

func (hc *HpaController) update_hpa_pod() {

	hpas, err := apirequest.GetAllHpas()
	if err != nil {
		return
	}
	fmt.Printf("hpas num:%d\n", len(hpas))

	hpaMap := make(map[string]string)
	for _, hpa := range hpas {
		value := hpa.MetaData.Namespace + "/" + hpa.MetaData.Name
		hpaMap[hpa.MetaData.UID] = value
	}

	pods, err := apirequest.GetAllPods()
	if err != nil {
		return
	}
	for _, pod := range pods {
		if pod.MetaData.Labels["hpa_uid"] == "" {
			continue
		}
		if _, ok := hpaMap[pod.MetaData.Labels["hpa_uid"]]; !ok {
			fmt.Println("hpa delete pod:", pod.MetaData.Name)
			hc.HpaDeletePod([]apiobj.Pod{pod}, 1)
		}
	}

	for _, hpa := range hpas {
		var num = 0
		var hpa_pods []apiobj.Pod
		for _, pod := range pods {

			for key, value := range hpa.Spec.Selector.MatchLabels {
				if pod.MetaData.Labels[key] == value {
					if pod.MetaData.Labels["hpa_uid"] == "" {
						pod.MetaData.Labels["hpa_uid"] = hpa.MetaData.UID

						url := apiconfig.URL_Pod
						url = apiconfig.GetApiServerUrl() + url
						url = strings.Replace(url, ":namespace", pod.MetaData.Namespace, -1)
						url = strings.Replace(url, ":name", pod.MetaData.Name, -1)

						apirequest.PutRequest(url, &pod)
					}

					num++
					hpa_pods = append(hpa_pods, pod)
				}
			}
		}
		fmt.Printf("hpa exist pod num:%d\n", len(hpa_pods))

		if len(hpa_pods) == 0 {
			continue
		}
		podCpuUsage := hc.getPodCpuUsage(hpa_pods)
		podMemUsage := hc.getPodMemUsage(hpa_pods)
		fmt.Printf("hpa pod cpu usage:%f\n", podCpuUsage)
		fmt.Printf("hpa mem usage:%f\n", podMemUsage)

		targetReplicas := hc.getTargetReplicas(hpa, podCpuUsage, podMemUsage)
		fmt.Printf("hpa target replicas:%d\n", targetReplicas)

		if targetReplicas > num {
			fmt.Printf("hpa add pod num:%d\n", targetReplicas-num)
			hc.HpaAddPod(hpa_pods[0], targetReplicas-num, hpa.MetaData)
		} else if targetReplicas < num {
			fmt.Printf("hpa delete pod num:%d\n", num-targetReplicas)
			hc.HpaDeletePod(hpa_pods, num-targetReplicas)
		}
		hpa.Status.Replicas = targetReplicas
		hpa.Status.CpuUsage = podCpuUsage * float64(targetReplicas)
		hpa.Status.MemUsage = podMemUsage * float64(targetReplicas)
		updateHpaStatus(hpa)
	}

}

func updateHpaStatus(hpa apiobj.Hpa) {

	URL := apiconfig.URL_Hpa
	URL = strings.Replace(URL, ":namespace", hpa.MetaData.Namespace, -1)
	URL = strings.Replace(URL, ":name", hpa.MetaData.Name, -1)
	HttpUrl := apiconfig.GetApiServerUrl() + URL

	jsonData, err := json.Marshal(hpa)
	if err != nil {
		fmt.Println("marshal hpa error")
		return
	}
	req, err := http.NewRequest(http.MethodPut, HttpUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("create put request error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("put error:", err)
		return
	}
	defer response.Body.Close()
}

func (hc *HpaController) HpaAddPod(podTemplate apiobj.Pod, num int, hpaMeta apiobj.MetaData) {

	pod := apiobj.Pod{
		ApiVersion: "v1",
		Kind:       "Pod",
		MetaData:   podTemplate.MetaData,
		Spec:       podTemplate.Spec,
	}

	pod.MetaData.Labels["hpa_uid"] = hpaMeta.UID
	oldPodName := pod.MetaData.Name
	oldContainerName := []string{}
	for _, container := range pod.Spec.Containers {
		oldContainerName = append(oldContainerName, container.Name)
	}

	url := apiconfig.URL_Pod
	url = apiconfig.GetApiServerUrl() + url
	for i := 0; i < num; i++ {
		pod.MetaData.Name = oldPodName + "-" + uuid.New().String()
		for id := range oldContainerName {
			pod.Spec.Containers[id].Name = oldContainerName[id] + "-" + uuid.New().String()
		}

		if pod.MetaData.Namespace == "" {
			pod.MetaData.Namespace = "default"
		}

		url = strings.Replace(url, ":namespace", pod.MetaData.Namespace, -1)
		url = strings.Replace(url, ":name", pod.MetaData.Name, -1)

		apirequest.PostRequest(url, &pod)
	}

}

func (hc *HpaController) HpaDeletePod(existPods []apiobj.Pod, num int) {

	url := apiconfig.URL_Pod
	url = apiconfig.GetApiServerUrl() + url
	for i := 0; i < num; i++ {
		pod := existPods[i]
		url = strings.Replace(url, ":namespace", pod.MetaData.Namespace, -1)
		url = strings.Replace(url, ":name", pod.MetaData.Name, -1)
		apirequest.DeleteRequest(url)
	}
}

func (hc *HpaController) getTargetReplicas(hpa apiobj.Hpa, podCpuUsage float64, podMemUsage float64) int {

	var targetReplicas = 0

	var targetCpuPercent = hpa.Spec.Metrics.CpuMetric.Target
	var targetMemPercent = hpa.Spec.Metrics.MemMetric.Target

	var cpuMaxReplicas = int(targetCpuPercent / float64(podCpuUsage))
	var memMaxReplicas = int(targetMemPercent / float64(podMemUsage))

	if cpuMaxReplicas < memMaxReplicas {
		targetReplicas = cpuMaxReplicas
	} else {
		targetReplicas = memMaxReplicas
	}

	if targetReplicas < hpa.Spec.MinReplicas {
		targetReplicas = hpa.Spec.MinReplicas
	} else if targetReplicas > hpa.Spec.MaxReplicas {
		targetReplicas = hpa.Spec.MaxReplicas
	}

	return targetReplicas
}

func (hc *HpaController) getPodCpuUsage(pods []apiobj.Pod) float64 {
	var cpuUsage = 0.0

	for _, pod := range pods {
		for _, container := range pod.Spec.Containers {
			cpuUsage += container.Resources.CPU
		}
	}

	cpuUsage = cpuUsage / float64(len(pods))
	return cpuUsage
}

func (hc *HpaController) getPodMemUsage(pods []apiobj.Pod) float64 {
	var memoryUsage = 0.0

	for _, pod := range pods {
		for _, container := range pod.Spec.Containers {
			memoryUsage += container.Resources.Memory
		}
	}

	memoryUsage = memoryUsage / float64(len(pods))
	return memoryUsage
}
