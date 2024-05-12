package controllers

import (
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/apirequest"
	"minik8s/pkg/config/apiconfig"
	"minik8s/tools/runner"
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
	rr.RunLoop(10*time.Second, 10*time.Second, hc.update_hpa_pod)
}

func (hc *HpaController) update_hpa_pod() {

	hpas, err := apirequest.GetAllHpas()
	if err != nil {
		return
	}
	fmt.Printf("hpas num:%d\n", len(hpas))
	if len(hpas) == 0 {
		return
	}

	hpaMap := make(map[string]string)
	for _, hpa := range hpas {
		value := hpa.MetaData.Namespace + "/" + hpa.MetaData.Name
		hpaMap[hpa.MetaData.UID] = value
	}

	pods, err := apirequest.GetAllPods()
	if err != nil {
		return
	}
	if len(pods) == 0 {
		return
	}
	for _, pod := range pods {
		if pod.MetaData.Labels["hpa_uid"] == "" {
			continue
		}
		if _, ok := hpaMap[pod.MetaData.Labels["hpa_uid"]]; !ok {
			fmt.Println("delete pod:", pod.MetaData.Name)
		}
	}

	for _, hpa := range hpas {
		var num = 0
		var hpa_pods []apiobj.Pod
		for _, pod := range pods {

			for key, value := range hpa.Spec.Selector.MatchLabels {
				if pod.MetaData.Labels[key] == value {
					num++
					hpa_pods = append(hpa_pods, pod)
				}
			}
		}
		fmt.Printf("existing hpa pod num:%d\n", len(hpa_pods))
		if len(hpa_pods) == 0 {
			continue
		}

		if hpa.Spec.MinReplicas > num {
			hc.HpaAddPod(hpa_pods[0], hpa.Spec.MinReplicas-num,hpa.MetaData)
		}
		if hpa.Spec.MaxReplicas < num {
			hc.HpaDeletePod(hpa_pods, num-hpa.Spec.MaxReplicas)
		}

	}

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
	oldContainerName :=  []string{}
	for _, container := range pod.Spec.Containers {
		oldContainerName = append(oldContainerName, container.Name)
	}

	for i := 0; i < num; i++ {
		pod.MetaData.Name = oldPodName + "-" + uuid.New().String()
		for id := range oldContainerName {
			pod.Spec.Containers[id].Name = oldContainerName[id] + "-" + uuid.New().String()
		}

		url := apiconfig.URL_Pod
		if pod.MetaData.Namespace == "" {
			pod.MetaData.Namespace = "default"
		}


		






}

func (hc *HpaController) HpaDeletePod(existPods []apiobj.Pod, num int) {

}
