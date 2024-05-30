package handler

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/etcd"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetGlobalPVCs(c *gin.Context) {
	fmt.Println("getGlobalPVCs")
	key := etcd.PATH_EtcdPVCs
	var resList []string
	resList, err := etcd.EtcdKV.GetPrefix(key)
	if err != nil {
		c.JSON(500, gin.H{"get": "fail"})
	}
	c.JSON(200, gin.H{
		"data": resList,
	})

}

func GetAllPVCs(c *gin.Context) {
	fmt.Println("getAllPVCs")
	namespace := c.Param("namespace")
	key := fmt.Sprintf(etcd.PATH_EtcdPVCs+"/%s", namespace)
	var resList []string
	resList, err := etcd.EtcdKV.GetPrefix(key)
	if err != nil {
		c.JSON(500, gin.H{"get": "fail"})
	}
	c.JSON(200, gin.H{
		"data": resList,
	})
}

func AddPVC(c *gin.Context) {
	fmt.Println("addPVC")
	var pvc apiobj.PVC
	c.ShouldBind(&pvc)
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdPVCs+"/%s/%s", namespace, name)

	pvc.MetaData.UID = uuid.New().String()[:16]

	pvcJson, err := json.Marshal(pvc)
	if err != nil {
		c.JSON(500, gin.H{"add": "fail"})
	}

	etcd.EtcdKV.Put(key, pvcJson)
	c.JSON(200, gin.H{"add": string(pvcJson)})
}

func DeletePVC(c *gin.Context) {
	fmt.Println("deletePVC")
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdPVCs+"/%s/%s", namespace, name)

	err := etcd.EtcdKV.Delete(key)
	if err != nil {
		c.JSON(500, gin.H{"delete": "fail"})
	}
	c.JSON(200, gin.H{"delete": "success"})
}

func UpdatePVC(c *gin.Context) {
	fmt.Println("updatePVC")

	var pvc apiobj.PVC
	c.ShouldBind(&pvc)
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdPVCs+"/%s/%s", namespace, name)

	pvcJson, err := json.Marshal(pvc)
	if err != nil {
		c.JSON(500, gin.H{"update": "fail"})
	}
	err = etcd.EtcdKV.Put(key, pvcJson)
	if err != nil {
		c.JSON(500, gin.H{"update": "fail"})
	}
}

func GetPVC(c *gin.Context) {
	fmt.Println("getPVC")
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdPVCs+"/%s/%s", namespace, name)

	res, err := etcd.EtcdKV.Get(key)
	if err != nil {
		c.JSON(500, gin.H{"get": "fail"})
	}
	c.JSON(200, gin.H{
		"data": string(res),
	})
}