package handler

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/etcd"
	"minik8s/pkg/message"
	"net/http"

	"github.com/gin-gonic/gin"
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
}

func DeletePod(c *gin.Context) {
	fmt.Println("deletePod")
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdPods+"/%s/%s", namespace, name)
	err := etcd.EtcdKV.Delete(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"delete": "fail"})
	}
	c.JSON(http.StatusOK, gin.H{"delete": "success"})
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
