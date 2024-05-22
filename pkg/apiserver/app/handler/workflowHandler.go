package handler

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/etcd"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetGlobalWorkflows(c *gin.Context) {
	fmt.Println("getGlobalWorkflows")
	key := etcd.PATH_EtcdWorkflows
	var resList []string
	resList, err := etcd.EtcdKV.GetPrefix(key)
	if err != nil {
		c.JSON(500, gin.H{"get": "fail"})
	}
	c.JSON(200, gin.H{
		"data": resList,
	})
}

func GetAllWorkflows(c *gin.Context) {
	fmt.Println("getAllWorkflows")
	namespace := c.Param("namespace")
	key := fmt.Sprintf(etcd.PATH_EtcdWorkflows+"/%s", namespace)
	var resList []string
	resList, err := etcd.EtcdKV.GetPrefix(key)
	if err != nil {
		c.JSON(500, gin.H{"get": "fail"})
	}
	c.JSON(200, gin.H{
		"data": resList,
	})
}

func AddWorkflow(c *gin.Context) {
	fmt.Println("addWorkflow")
	var workflow apiobj.Workflow
	c.ShouldBind(&workflow)
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdWorkflows+"/%s/%s", namespace, name)

	workflow.MetaData.UID = uuid.New().String()

	workflowJson, err := json.Marshal(workflow)
	if err != nil {
		c.JSON(500, gin.H{"add": "fail"})
	}

	etcd.EtcdKV.Put(key, workflowJson)
	c.JSON(200, gin.H{"add": string(workflowJson)})
}

func DeleteWorkflow(c *gin.Context) {
	fmt.Println("deleteWorkflow")
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdWorkflows+"/%s/%s", namespace, name)

	etcd.EtcdKV.Delete(key)
	c.JSON(200, gin.H{"delete": "success"})
}

func UpdateWorkflow(c *gin.Context) {
	fmt.Println("updateWorkflow")
	var workflow apiobj.Workflow
	c.ShouldBind(&workflow)
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdWorkflows+"/%s/%s", namespace, name)

	workflowJson, err := json.Marshal(workflow)
	if err != nil {
		c.JSON(500, gin.H{"update": "fail"})
	}

	etcd.EtcdKV.Put(key, workflowJson)
	c.JSON(200, gin.H{"update": string(workflowJson)})
}

func GetWorkflow(c *gin.Context) {
	fmt.Println("getWorkflow")
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdWorkflows+"/%s/%s", namespace, name)
	res, err := etcd.EtcdKV.Get(key)
	if err != nil {
		c.JSON(500, gin.H{"get": "fail"})
	}
	c.JSON(200, gin.H{
		"data": string(res),
	})
}