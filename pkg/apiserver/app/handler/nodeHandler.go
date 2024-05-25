package handler

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	etcd "minik8s/pkg/etcd"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetNodes(c *gin.Context) {
	fmt.Println("getNodes")
	key := etcd.PATH_EtcdNodes
	res, err := etcd.EtcdKV.GetPrefix(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"get": "fail"})
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

func AddNode(c *gin.Context) {
	fmt.Println("addNode")

	var node apiobj.Node
	c.ShouldBind(&node)
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdNodes+"/%s", name)
	nodeJson, err := json.Marshal(node)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"add": "fail"})
	}
	etcd.EtcdKV.Put(key, nodeJson)

	c.JSON(http.StatusOK, gin.H{"add": "success"})
}

func DeleteNode(c *gin.Context) {
	fmt.Println("deleteNode")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdNodes+"/%s", name)
	err := etcd.EtcdKV.Delete(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"delete": "fail"})
	}
	c.JSON(http.StatusOK, gin.H{"delete": "success"})
}

func UpdateNode(c *gin.Context) {
	fmt.Println("updateNode")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdNodes+"/%s", name)
	value := []byte("node")
	err := etcd.EtcdKV.Put(key, value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"update": "fail"})
	}
	c.JSON(http.StatusOK, gin.H{"update": "success"})
}

func GetNode(c *gin.Context) {

	fmt.Print("getNode")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdNodes+"/%s", name)
	res, err := etcd.EtcdKV.Get(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"get": "fail"})
	}
	c.JSON(http.StatusOK, gin.H{"data": string(res)})
}

func GetNodePods(c *gin.Context) {
	fmt.Println("getNodePods")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdPods+"%s/pods", name)
	res, err := etcd.EtcdKV.Get(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"get": "fail"})
	}
	c.JSON(http.StatusOK, gin.H{"data": string(res)})
}

func GetNodeStatus(c *gin.Context) {
	fmt.Println("getNodeStatus")
	fmt.Println("getNodePods")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdPods+"%s/status", name)
	res, err := etcd.EtcdKV.Get(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"get": "fail"})
	}
	c.JSON(http.StatusOK, gin.H{"data": string(res)})
}
