package controllers

import (
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/apirequest"
	"minik8s/pkg/config/apiconfig"
	"minik8s/tools/runner"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ReplicaController struct {
	replicaMap map[string]string
}

func NewReplicaController() *ReplicaController {
	return &ReplicaController{
		replicaMap: make(map[string]string),
	}
}

func (rc *ReplicaController) Run() {
	rr := runner.NewRunner()
 	rr.RunLoop(5*time.Second, 5*time.Second, rc.update_pod_num)
}

func (rc *ReplicaController) update_pod_num() {
	replicasets, err := apirequest.GetAllReplicaSets()
	if err != nil {
		return
	}

	fmt.Printf("replicasets num:%d\n", len(replicasets))
	
	replicaMap := make(map[string]string)
	for _, rs := range replicasets {
		value := rs.MetaData.Namespace + "/" + rs.MetaData.Name
		replicaMap[rs.MetaData.UID] = value
		fmt.Printf("replicasets uid :%s  set\n", rs.MetaData.UID)
	}

	pods, err := apirequest.GetAllPods()
	if err != nil {
		return
	}
	for _, pod := range pods {
		if pod.MetaData.Labels["replica_uid"] == "" {
			continue
		}
		if _, ok := replicaMap[pod.MetaData.Labels["replica_uid"]]; !ok {
			fmt.Print("delete pod\n")
			rc.DeleteReplica([]apiobj.Pod{pod}, 1)
		}
	}

	for _, replicaset := range replicasets {
		var num = 0
		var replica_pods []apiobj.Pod
		for _, pod := range pods {

			for key, value := range replicaset.Spec.Selector.MatchLabels {
				fmt.Printf("key:%s value:%s pod value:%s\n", key, value,pod.MetaData.Labels[key])
				if pod.MetaData.Labels[key] == value {
					num++
					replica_pods = append(replica_pods, pod)
				}
			}
		}

		fmt.Printf("existing replica pod num:%d\n",num)

		if num < replicaset.Spec.Replicas {
			fmt.Printf("add pod num:%d\n",replicaset.Spec.Replicas - num)
			rc.AddReplica(replicaset.Spec.Template, replicaset.Spec.Replicas - num, replicaset.MetaData)
		} else if num > replicaset.Spec.Replicas {
			fmt.Printf("delete pod num:%d\n",num - replicaset.Spec.Replicas)
			rc.DeleteReplica(replica_pods, num - replicaset.Spec.Replicas)
		}
	}

}
	
func (rc *ReplicaController) AddReplica(podTemplate apiobj.PodTemplateSpec, num int, replicaMeta apiobj.MetaData) error {
	pod := apiobj.Pod{
		ApiVersion: "v1",
		Kind:       "Pod",
		MetaData:   podTemplate.MetaData,
		Spec:       podTemplate.Spec,
	}
	pod.MetaData.Labels["replica_uid"] = replicaMeta.UID

	oldPodName := pod.MetaData.Name
	oldContainerName := []string{}
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
		url = strings.Replace(url, ":namespace", pod.MetaData.Namespace, -1)
		url = strings.Replace(url, ":name", pod.MetaData.Name, -1)
		url = apiconfig.GetApiServerUrl() + url

		apirequest.PostRequest(url, &pod)

	}

	return nil
}

func (rc *ReplicaController) DeleteReplica(existPods []apiobj.Pod, num int) error {
	
	url := apiconfig.URL_Pod
	url = apiconfig.GetApiServerUrl() + url
	for i := 0; i < num; i++ {

		pod := existPods[i]

		if pod.MetaData.Namespace == "" {
			pod.MetaData.Namespace = "default"
		}
		url = strings.Replace(url, ":namespace", pod.MetaData.Namespace, -1)
		url = strings.Replace(url, ":name", pod.MetaData.Name, -1)
		apirequest.DeleteRequest(url)
	}
	return nil
}
