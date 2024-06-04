package handler

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/etcd"
	"minik8s/pkg/message"
	nginxmanager "minik8s/pkg/nginx/app"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetGlobalPods(c *gin.Context) {
	fmt.Println("getGlobalPods")
	key := etcd.PATH_EtcdPods
	var resList []string
	resList, err := etcd.EtcdKV.GetPrefix(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"get": "fail"})
	}
	c.JSON(http.StatusOK, gin.H{
		"data": resList,
	})

}

func GetAllPods(c *gin.Context) {
	fmt.Println("getAllPods")
	namespace := c.Param("namespace")
	key := fmt.Sprintf(etcd.PATH_EtcdPods+"/%s", namespace)
	var resList []string
	resList, err := etcd.EtcdKV.GetPrefix(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"get": "fail"})
	}
	c.JSON(http.StatusOK, gin.H{
		"data": resList,
	})

}

func AddPod(c *gin.Context) {
	fmt.Println("addPod")
	var pod apiobj.Pod
	c.ShouldBind(&pod)
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdPods+"/%s/%s", namespace, name)
	pod.MetaData.UID = uuid.New().String()[:16]

	//pv pvc handle
	if len(pod.Spec.Volumes) > 0 && pod.Spec.Volumes[0].HostPath.Path == "" {

		if pod.Spec.Volumes[0].PersistentVolumeClaim.ClaimName != "" {
			pvcKey := fmt.Sprintf(etcd.PATH_EtcdPVCs+"/%s/%s", namespace, pod.Spec.Volumes[0].PersistentVolumeClaim.ClaimName)
			pvcRes, _ := etcd.EtcdKV.Get(pvcKey)
			var pvc apiobj.PVC
			json.Unmarshal([]byte(pvcRes), &pvc)

			pvKey := fmt.Sprintf(etcd.PATH_EtcdPVs+"/%s", namespace)
			pvList, err := etcd.EtcdKV.GetPrefix(pvKey)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"get": "fail"})
			}
			var pvs []apiobj.PV
			for _, item := range pvList {
				var pv apiobj.PV
				json.Unmarshal([]byte(item), &pv)
				pvs = append(pvs, pv)
			}

			var pvExist = false
			for _, pv := range pvs {
				if pv.Spec.StorageClassName == pvc.Spec.StorageClassName {
					mntPath := apiobj.NfsMntPath
					newPath := fmt.Sprintf("%s/%s", mntPath, pv.MetaData.Name)
					pod.Spec.Volumes[0].HostPath.Path = newPath
					pvExist = true
					break
				}
			}

			if !pvExist {
				pv := apiobj.PV{
					APIVersion: "v1",
					Kind:       "PV",
					MetaData: apiobj.MetaData{
						Name:      pod.MetaData.Name + "-pv",
						Namespace: pod.MetaData.Namespace,
						UID:       uuid.New().String()[:16],
					},
					Spec: apiobj.PVSpec{
						StorageClassName: pvc.Spec.StorageClassName,
					},
				}

				mntPath := apiobj.NfsMntPath
				newPath := fmt.Sprintf("%s/%s", mntPath, pv.MetaData.Name)

				err = os.Mkdir(newPath, 0755)
				if err != nil {
					fmt.Println(err)
				}
				pod.Spec.Volumes[0].HostPath.Path = newPath

				pvKey := fmt.Sprintf(etcd.PATH_EtcdPVs+"/%s/%s", pv.MetaData.Namespace, pv.MetaData.Name)
				pvJson, _ := json.Marshal(pv)
				etcd.EtcdKV.Put(pvKey, pvJson)

			}

		}

	}

	podJson, err := json.Marshal(pod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"add": "fail"})
	}

	etcd.EtcdKV.Put(key, podJson)
	c.JSON(http.StatusOK, gin.H{"add": string(podJson)})

	msg := message.Message{
		Type:    "Add",
		URL:     key,
		Name:    name,
		Content: string(podJson),
	}
	msgJson, _ := json.Marshal(msg)
	p := message.NewPublisher()
	defer p.Close()
	p.Publish(message.ScheduleQueue, msgJson)

	//replicaset handle

}

func DeletePod(c *gin.Context) {
	fmt.Println("deletePod")
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdPods+"/%s/%s", namespace, name)

	res, _ := etcd.EtcdKV.Get(key)
	var pod apiobj.Pod
	json.Unmarshal([]byte(res), &pod)

	err := etcd.EtcdKV.Delete(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"delete": "fail"})
	}
	c.JSON(http.StatusOK, gin.H{"delete": "success"})

	msg := message.Message{
		Type:    "Delete",
		URL:     key,
		Name:    name,
		Content: string(res),
	}
	msgJson, _ := json.Marshal(msg)
	p := message.NewPublisher()
	defer p.Close()

	podQue := fmt.Sprintf(message.PodQueue+"-%s", pod.Status.NodeName)
	p.Publish(podQue, msgJson)

	// service handle
	if pod.MetaData.Labels["app"] != "" {
		svcKey := fmt.Sprintf(etcd.PATH_EtcdServices+"/%s/%s", namespace, pod.MetaData.Labels["app"])
		res, _ := etcd.EtcdKV.Get(svcKey)
		var svc apiobj.Service
		json.Unmarshal([]byte(res), &svc)

		var svcports []string
		var podports []string
		for _, port := range svc.Spec.Ports {
			svcports = append(svcports, fmt.Sprintf("%d", port.Port))
			podports = append(podports, fmt.Sprintf("%d", port.TargetPort))
			// update nginx config
			nginxmanager.DeleteServiceRule(svc.Spec.ClusterIP, uint16(port.Port), pod.Status.PodIP, uint16(port.TargetPort))
		}

		svcMsg := apiobj.PodSvcMsg{
			SvcIp:    svc.Spec.ClusterIP,
			SvcPorts: svcports,
			PodIp:    pod.Status.PodIP,
			PodPorts: podports,
		}

		svcMsgJson, _ := json.Marshal(svcMsg)
		msg := message.Message{
			Type:    "Update",
			URL:     svcKey,
			Name:    "Delete",
			Content: string(svcMsgJson),
		}
		msgJson, _ := json.Marshal(msg)
		p := message.NewPublisher()
		defer p.Close()
		
		nodeKey := etcd.PATH_EtcdNodes
		resList, _ := etcd.EtcdKV.GetPrefix(nodeKey)

		for _, item := range resList {
			var node apiobj.Node
			json.Unmarshal([]byte(item), &node)
			que := fmt.Sprintf(message.ServiceQueue+"-%s", node.MetaData.Name)
			p.Publish(que, msgJson)
		}

	}
}

func UpdatePod(c *gin.Context) {
	fmt.Println("updatePod")

	var pod apiobj.Pod
	c.ShouldBind(&pod)
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdPods+"/%s/%s", namespace, name)

	podJson, err := json.Marshal(pod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"update": "fail"})
	}

	etcd.EtcdKV.Put(key, podJson)

	//service handle
	if pod.MetaData.Labels["svc"] != "" && pod.Status.PodIP != "" {
		svcKey := fmt.Sprintf(etcd.PATH_EtcdServices+"/%s/%s", namespace, pod.MetaData.Labels["svc"])
		res, err := etcd.EtcdKV.Get(svcKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"get": "fail"})
			fmt.Println("no service")
			return
		}
		var svc apiobj.Service
		json.Unmarshal([]byte(res), &svc)

		var svcports []string
		var podports []string
		for _, port := range svc.Spec.Ports {
			svcports = append(svcports, fmt.Sprintf("%d", port.Port))
			podports = append(podports, fmt.Sprintf("%d", port.TargetPort))
			// update nginx config
			nginxmanager.AddServiceRule(svc.Spec.ClusterIP, uint16(port.Port), pod.Status.PodIP, uint16(port.TargetPort))
		}

		svcMsg := apiobj.PodSvcMsg{
			SvcIp:    svc.Spec.ClusterIP,
			SvcPorts: svcports,
			PodIp:    pod.Status.PodIP,
			PodPorts: podports,
		}

		svcMsgJson, _ := json.Marshal(svcMsg)
		msg := message.Message{
			Type:    "Update",
			URL:     key,
			Name:    "Add",
			Content: string(svcMsgJson),
		}
		msgJson, _ := json.Marshal(msg)

		p := message.NewPublisher()
		defer p.Close()

		nodeKey := etcd.PATH_EtcdNodes
		resList, _ := etcd.EtcdKV.GetPrefix(nodeKey)

		for _, item := range resList {
			var node apiobj.Node
			json.Unmarshal([]byte(item), &node)
			que := fmt.Sprintf(message.ServiceQueue+"-%s", node.MetaData.Name)
			p.Publish(que, msgJson)
		}

	}

	c.JSON(http.StatusOK, gin.H{"update": string(podJson)})
}

func GetPod(c *gin.Context) {
	fmt.Println("getPod")
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdPods+"/%s/%s", namespace, name)
	res, err := etcd.EtcdKV.Get(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"get": "fail"})
	}
	c.JSON(http.StatusOK, gin.H{
		"data": string(res),
	})
}

func GetPodStatus(c *gin.Context) {
	fmt.Println("getPodStatus")
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdPods+"/%s/%s/status", namespace, name)
	res, err := etcd.EtcdKV.Get(key)
	var pod apiobj.Pod
	json.Unmarshal([]byte(res), &pod)

	var status = pod.Status
	statusJson, _ := json.Marshal(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"get": "fail"})
	}
	c.JSON(http.StatusOK, gin.H{
		"data": string(statusJson),
	})
}

func UpdatePodStatus(c *gin.Context) {
	fmt.Println("updatePodStatus")

	var podStatus apiobj.PodStatus
	c.ShouldBind(&podStatus)
	namespace := c.Param("namespace")
	name := c.Param("name")

	key := fmt.Sprintf(etcd.PATH_EtcdPods+"/%s/%s", namespace, name)
	res, err := etcd.EtcdKV.Get(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"get": "fail"})
	}

	var pod apiobj.Pod
	json.Unmarshal([]byte(res), &pod)
	pod.Status = podStatus

	podJson, _ := json.Marshal(pod)
	etcd.EtcdKV.Put(key, podJson)
	c.JSON(http.StatusOK, gin.H{"update": string(podJson)})
}
