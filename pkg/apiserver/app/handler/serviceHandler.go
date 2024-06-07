package handler

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/config/serviceconfig"
	"minik8s/pkg/etcd"
	"minik8s/pkg/message"
	"net/http"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

func GetGlobalServices(c *gin.Context) {
	// get global services
	key := etcd.PATH_EtcdServices
	var resList []string
	resList, err := etcd.EtcdKV.GetPrefix(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"get": "fail"})
	}
	c.JSON(http.StatusOK, gin.H{
		"data": resList,
	})
}

func GetAllServices(c *gin.Context) {
	// get all services
	namespace := c.Param("namespace")
	key := fmt.Sprintf(etcd.PATH_EtcdServices+"/%s", namespace)
	var resList []string
	resList, err := etcd.EtcdKV.GetPrefix(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"get": "fail"})
	}
	c.JSON(http.StatusOK, gin.H{
		"data": resList,
	})
}

func GetService(c *gin.Context) {
	// get service
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdServices+"/%s/%s", namespace, name)
	res, err := etcd.EtcdKV.Get(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"get": "fail"})
	}
	c.JSON(http.StatusOK, gin.H{
		"data": string(res),
	})
}

func AddService(c *gin.Context) {
	// create service
	var service apiobj.Service
	c.ShouldBind(&service)
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdServices+"/%s/%s", namespace, name)

	service.MetaData.UID = uuid.New().String()[:16]

	service.Spec.ClusterIP = serviceconfig.AllocateIp()

	// update nginx config
	// nginxmanager.AddServiceIPVS(service.Spec)
	serviceJson, err := json.Marshal(service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"add": "fail"})
	}

	etcd.EtcdKV.Put(key, serviceJson)
	c.JSON(http.StatusOK, gin.H{"add": string(serviceJson)})

	msg := message.Message{
		Type:    "Add",
		URL:     key,
		Name:    name,
		Content: string(serviceJson),
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

func UpdateService(c *gin.Context) {
	// update service
	var service apiobj.Service
	c.ShouldBind(&service)
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdServices+"/%s/%s", namespace, name)

	serviceJson, err := json.Marshal(service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"update": "fail"})
	}

	etcd.EtcdKV.Put(key, serviceJson)
	c.JSON(http.StatusOK, gin.H{"update": string(serviceJson)})

	msg := message.Message{
		Type:    "Update",
		URL:     key,
		Name:    name,
		Content: string(serviceJson),
	}
	msgJson, _ := json.Marshal(msg)
	p := message.NewPublisher()
	defer p.Close()
	p.Publish(message.ScheduleQueue, msgJson)
}

func DeleteService(c *gin.Context) {
	// delete service
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdServices+"/%s/%s", namespace, name)

	var service apiobj.Service
	res, _ := etcd.EtcdKV.Get(key)
	json.Unmarshal([]byte(res), &service)
	serviceIp := service.Spec.ClusterIP
	serviceconfig.ReleaseIp(serviceIp)
	// update nginx config
	// nginxmanager.DeleteServiceIPVS(service.Spec)

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

	nodeKey := etcd.PATH_EtcdNodes
	resList, _ := etcd.EtcdKV.GetPrefix(nodeKey)

	for _, item := range resList {
		var node apiobj.Node
		json.Unmarshal([]byte(item), &node)
		que := fmt.Sprintf(message.ServiceQueue+"-%s", node.MetaData.Name)
		p.Publish(que, msgJson)
	}
}

func GetServiceStatus(c *gin.Context) {
	// get service status
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdServices+"/%s/%s/status", namespace, name)
	res, err := etcd.EtcdKV.Get(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"get": "fail"})
	}
	var service apiobj.Service
	json.Unmarshal([]byte(res), &service)

	var status = service.Status
	statusJson, _ := json.Marshal(status)
	c.JSON(http.StatusOK, gin.H{
		"data": string(statusJson),
	})
}

func UpdateServiceStatus(c *gin.Context) {
	// update service status
	var status apiobj.ServiceStatus
	c.ShouldBind(&status)
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdServices+"/%s/%s", namespace, name)

	res, err := etcd.EtcdKV.Get(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"update": "fail"})
	}
	var service apiobj.Service
	json.Unmarshal([]byte(res), &service)
	service.Status = status

	serviceJson, err := json.Marshal(service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"update": "fail"})
	}

	etcd.EtcdKV.Put(key, serviceJson)
	c.JSON(http.StatusOK, gin.H{"update": string(serviceJson)})
}
