package handler

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/etcd"
	"net/http"

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
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdServices+"/%s", name)
	res, err := etcd.EtcdKV.Get(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"get": "fail"})
	}
	c.JSON(http.StatusOK, gin.H{
		"data": res,
	})
}

func AddService(c *gin.Context) {
	// create service
	var service apiobj.Service
	c.ShouldBind(&service)
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdServices+"/%s/%s", namespace, name)

	serviceJson, err := json.Marshal(service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"add": "fail"})
	}

	etcd.EtcdKV.Put(key, serviceJson)
	c.JSON(http.StatusOK, gin.H{"add": string(serviceJson)})
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
}

func DeleteService(c *gin.Context) {
	// delete service
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdServices+"/%s/%s", namespace, name)
	err := etcd.EtcdKV.Delete(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"delete": "fail"})
	}
	c.JSON(http.StatusOK, gin.H{"delete": "success"})
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
