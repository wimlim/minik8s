package handler

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/etcd"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetGlobalFunctions(c *gin.Context) {
	fmt.Println("getGlobalFunctions")
	key := etcd.PATH_EtcdFunctions
	var resList []string
	resList, err := etcd.EtcdKV.GetPrefix(key)
	if err != nil {
		c.JSON(500, gin.H{"get": "fail"})
	}
	c.JSON(200, gin.H{
		"data": resList,
	})
}

func GetAllFunctions(c *gin.Context) {
	fmt.Println("getAllFunctions")
	namespace := c.Param("namespace")
	key := fmt.Sprintf(etcd.PATH_EtcdFunctions+"/%s", namespace)
	var resList []string
	resList, err := etcd.EtcdKV.GetPrefix(key)
	if err != nil {
		c.JSON(500, gin.H{"get": "fail"})
	}
	c.JSON(200, gin.H{
		"data": resList,
	})
}

func AddFunction(c *gin.Context) {
	fmt.Println("addFunction")
	var function apiobj.Function
	c.ShouldBind(&function)
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdFunctions+"/%s/%s", namespace, name)

	function.MetaData.UID = uuid.New().String()[:16]

	functionJson, err := json.Marshal(function)
	if err != nil {
		c.JSON(500, gin.H{"add": "fail"})
	}

	etcd.EtcdKV.Put(key, functionJson)
	c.JSON(200, gin.H{"add": string(functionJson)})
}

func DeleteFunction(c *gin.Context) {
	fmt.Println("deleteFunction")
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdFunctions+"/%s/%s", namespace, name)

	err := etcd.EtcdKV.Delete(key)
	if err != nil {
		c.JSON(500, gin.H{"delete": "fail"})
	}
	c.JSON(200, gin.H{"delete": "success"})
}

func UpdateFunction(c *gin.Context) {
	fmt.Println("updateFunction")
	var function apiobj.Function
	c.ShouldBind(&function)
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdFunctions+"/%s/%s", namespace, name)

	functionJson, err := json.Marshal(function)
	if err != nil {
		c.JSON(500, gin.H{"update": "fail"})
	}
	err = etcd.EtcdKV.Put(key, functionJson)
	if err != nil {
		c.JSON(500, gin.H{"update": "fail"})
	}
	c.JSON(200, gin.H{"update": string(functionJson)})
}

func GetFunction(c *gin.Context) {
	fmt.Println("getFunction")
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdFunctions+"/%s/%s", namespace, name)
	res, err := etcd.EtcdKV.Get(key)
	if err != nil {
		c.JSON(500, gin.H{"get": "fail"})
	}
	c.JSON(200, gin.H{
		"data": string(res),
	})
}
