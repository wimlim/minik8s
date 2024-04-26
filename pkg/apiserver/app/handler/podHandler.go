package handler

import (
	"fmt"
	"minik8s/pkg/etcd"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetGlobalPods(c *gin.Context){
	fmt.Println("getGlobalPods")
}

func GetAllPods(c *gin.Context){
	fmt.Println("getAllPods")
}

func AddPod(c *gin.Context){
	fmt.Println("addPod")
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdPods+"%s/%s", namespace, name)
	value := []byte("pod")
	err := etcd.EtcdKV.Put(key, value)
	if(err != nil){}
	c.JSON(http.StatusOK, gin.H{"add": "success"})
	
}

func DeletePod(c *gin.Context){
	fmt.Println("deletePod")
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdPods+"%s/%s", namespace, name)
	err := etcd.EtcdKV.Delete(key)	
	if(err != nil){}
	c.JSON(http.StatusOK, gin.H{"delete": "success"})
}

func UpdatePod(c *gin.Context){
	fmt.Println("updatePod")
	fmt.Println("addPod")
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdPods+"%s/%s", namespace, name)
	value := []byte("pod")
	err := etcd.EtcdKV.Put(key, value)
	if(err != nil){}
	c.JSON(http.StatusOK, gin.H{"update": "success"})
}

func GetPod(c *gin.Context){
	fmt.Println("getPod")
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdPods+"%s/%s", namespace, name)
	res,err := etcd.EtcdKV.Get(key);
	if(err != nil){}
	c.JSON(http.StatusOK, gin.H{
		"data": res,
	})
}

func GetPodStatus(c *gin.Context){
	fmt.Println("getPodStatus")
}