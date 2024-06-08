package handler

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/etcd"

	"github.com/gin-gonic/gin"
)

func GetGlobalJobs(c *gin.Context) {
	fmt.Println("getGlobalJobs")
	key := etcd.PATH_EtcdJobs
	var resList []string
	resList, err := etcd.EtcdKV.GetPrefix(key)
	if err != nil {
		c.JSON(500, gin.H{"get": "fail"})
	}
	c.JSON(200, gin.H{
		"data": resList,
	})

}

func GetAllJobs(c *gin.Context) {
	fmt.Println("getAllJobs")
	namespace := c.Param("namespace")
	key := fmt.Sprintf(etcd.PATH_EtcdJobs+"/%s", namespace)
	var resList []string
	resList, err := etcd.EtcdKV.GetPrefix(key)
	if err != nil {
		c.JSON(500, gin.H{"get": "fail"})
	}
	c.JSON(200, gin.H{
		"data": resList,
	})
}

func AddJob(c *gin.Context) {
	fmt.Println("addJob")
	var job apiobj.Job
	c.ShouldBind(&job)
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdJobs+"/%s/%s", namespace, name)

	jobJson, err := json.Marshal(job)
	if err != nil {
		c.JSON(500, gin.H{"add": "fail"})
	}

	etcd.EtcdKV.Put(key, jobJson)
	c.JSON(200, gin.H{"add": string(jobJson)})
}

func DeleteJob(c *gin.Context) {
	fmt.Println("deleteJob")
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdJobs+"/%s/%s", namespace, name)

	err := etcd.EtcdKV.Delete(key)
	if err != nil {
		c.JSON(500, gin.H{"delete": "fail"})
	}
	c.JSON(200, gin.H{"delete": "success"})
}

func UpdateJob(c *gin.Context) {
	fmt.Println("updateJob")

	var job apiobj.Job
	c.ShouldBind(&job)
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdJobs+"/%s/%s", namespace, name)

	jobJson, err := json.Marshal(job)
	if err != nil {
		c.JSON(500, gin.H{"update": "fail"})
	}
	err = etcd.EtcdKV.Put(key, jobJson)
	if err != nil {
		c.JSON(500, gin.H{"update": "fail"})
	}
	c.JSON(200, gin.H{"update": "success"})
}

func GetJob(c *gin.Context) {
	fmt.Println("getJob")
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdJobs+"/%s/%s", namespace, name)
	res, err := etcd.EtcdKV.Get(key)
	if err != nil {
		c.JSON(500, gin.H{"get": "fail"})
	}
	c.JSON(200, gin.H{"data": string(res)})
}