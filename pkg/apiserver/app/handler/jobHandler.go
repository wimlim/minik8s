package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/etcd"
	"minik8s/pkg/message"
	"os"

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

	job.Status.Phase = apiobj.Running
	jobJson, err := json.Marshal(job)
	if err != nil {
		c.JSON(500, gin.H{"add": "fail"})
	}

	etcd.EtcdKV.Put(key, jobJson)
	c.JSON(200, gin.H{"add": string(jobJson)})

	msg := message.Message{
		Type:    "Add",
		URL:     key,
		Name:    key,
		Content: string(jobJson),
	}
	msgJson, _ := json.Marshal(msg)
	p := message.NewPublisher()
	defer p.Close()

	p.Publish(message.JobQueue, msgJson)

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

func UpdateJobStatus(c *gin.Context) {
	fmt.Println("updateJobStatus")
	namespace := c.Param("namespace")
	name := c.Param("name")

	var jobStatus apiobj.JobStatus
	c.ShouldBind(&jobStatus)

	key := fmt.Sprintf(etcd.PATH_EtcdJobs+"/%s/%s", namespace, name)

	res, _ := etcd.EtcdKV.Get(key)
	var job apiobj.Job
	json.Unmarshal([]byte(res), &job)
	job.Status = jobStatus

	filepath := fmt.Sprintf("/tmp/results/%s.out", job.MetaData.Name)
	fd, err := os.Open(filepath)
	if err != nil {
		fmt.Println("open file error")
		return
	}
	defer fd.Close()
	content, err := io.ReadAll(fd)
	if err != nil {
		fmt.Println("read file error")
		c.JSON(500, gin.H{"update": "fail"})
	}

	job.Status.Result = string(content)

	jobJson, _ := json.Marshal(job)
	etcd.EtcdKV.Put(key, jobJson)
	c.JSON(200, gin.H{"update": "success"})
}
