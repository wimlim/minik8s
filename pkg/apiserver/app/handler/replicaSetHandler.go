package handler

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/etcd"
	"net/http"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

// GetGlobalReplicaSets is a function.
func GetGlobalReplicaSets(c *gin.Context) {
	key := etcd.PATH_EtcdReplicas
	var resList []string
	resList, err := etcd.EtcdKV.GetPrefix(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"get": "fail"})
	}
	c.JSON(http.StatusOK, gin.H{
		"data": resList,
	})
}

// GetAllReplicaSets is a function.
func GetAllReplicaSets(c *gin.Context) {
	namespace := c.Param("namespace")
	key := fmt.Sprintf(etcd.PATH_EtcdReplicas+"/%s", namespace)
	var resList []string
	resList, err := etcd.EtcdKV.GetPrefix(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"get": "fail"})
	}
	c.JSON(http.StatusOK, gin.H{
		"data": resList,
	})
}

func AddReplicaSet(c *gin.Context) {
	var replicaSet apiobj.ReplicaSet
	c.ShouldBind(&replicaSet)
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdReplicas+"/%s/%s", namespace, name)

	replicaSet.MetaData.UID = uuid.New().String()[:16]

	replicaSetJson, err := json.Marshal(replicaSet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"add": "fail"})
	}

	etcd.EtcdKV.Put(key, replicaSetJson)
	c.JSON(http.StatusOK, gin.H{"add": string(replicaSetJson)})
}

// DeleteReplicaSet is a function.
func DeleteReplicaSet(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdReplicas+"/%s/%s", namespace, name)
	err := etcd.EtcdKV.Delete(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"delete": "fail"})
	}
	c.JSON(http.StatusOK, gin.H{"delete": "success"})
}

func UpdateReplicaSet(c *gin.Context) {
	var replicaSet apiobj.ReplicaSet
	c.ShouldBind(&replicaSet)
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdReplicas+"/%s/%s", namespace, name)

	replicaSetJson, err := json.Marshal(replicaSet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"update": "fail"})
	}

	etcd.EtcdKV.Put(key, replicaSetJson)
	c.JSON(http.StatusOK, gin.H{"update": string(replicaSetJson)})
}

// GetReplicaSet is a function.
func GetReplicaSet(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdReplicas+"/%s/%s", namespace, name)
	res, err := etcd.EtcdKV.Get(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"get": "fail"})
	}
	c.JSON(http.StatusOK, gin.H{
		"data": string(res),
	})
}

// GetReplicaSetStatus is a function.
func GetReplicaSetStatus(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdReplicas+"/%s/%s/status", namespace, name)
	res, err := etcd.EtcdKV.Get(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"get": "fail"})
	}
	c.JSON(http.StatusOK, gin.H{
		"data": string(res),
	})
}

func UpdateReplicaSetStatus(c *gin.Context) {
	var replicaSetStatus apiobj.ReplicaSetStatus
	c.ShouldBind(&replicaSetStatus)
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdReplicas+"/%s/%s", namespace, name)

	res, err := etcd.EtcdKV.Get(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"get": "fail"})
	}

	var replicaSet apiobj.ReplicaSet
	json.Unmarshal([]byte(res), &replicaSet)
	replicaSet.Status = replicaSetStatus
	
	replicaSetJson, err := json.Marshal(replicaSet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"update": "fail"})
	}

	etcd.EtcdKV.Put(key, replicaSetJson)
	c.JSON(http.StatusOK, gin.H{"update": string(replicaSetJson)})
}
