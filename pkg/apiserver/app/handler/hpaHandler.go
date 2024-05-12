package handler

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/etcd"

	"github.com/gin-gonic/gin"
)

func GetGlobalHpas(c *gin.Context) {
	fmt.Println("getGlobalHpas")
	key := etcd.PATH_EtcdHpas
	var resList []string
	resList, err := etcd.EtcdKV.GetPrefix(key)
	if err != nil {
		c.JSON(500, gin.H{"get": "fail"})
	}
	c.JSON(200, gin.H{
		"data": resList,
	})

}

func GetAllHpas(c *gin.Context) {
	fmt.Println("getAllHpas")
	namespace := c.Param("namespace")
	key := fmt.Sprintf(etcd.PATH_EtcdHpas+"/%s", namespace)
	var resList []string
	resList, err := etcd.EtcdKV.GetPrefix(key)
	if err != nil {
		c.JSON(500, gin.H{"get": "fail"})
	}
	c.JSON(200, gin.H{
		"data": resList,
	})
}

func AddHpa(c *gin.Context) {
	fmt.Println("addHpa")
	var hpa apiobj.Hpa
	c.ShouldBind(&hpa)
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdHpas+"/%s/%s", namespace, name)

	hpaJson, err := json.Marshal(hpa)
	if err != nil {
		c.JSON(500, gin.H{"add": "fail"})
	}

	etcd.EtcdKV.Put(key, hpaJson)
	c.JSON(200, gin.H{"add": string(hpaJson)})
}

func DeleteHpa(c *gin.Context) {
	fmt.Println("deleteHpa")
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdHpas+"/%s/%s", namespace, name)
	err := etcd.EtcdKV.Delete(key)
	if err != nil {
		c.JSON(500, gin.H{"delete": "fail"})
	}
	c.JSON(200, gin.H{"delete": "success"})
}

func UpdateHpa(c *gin.Context) {
	fmt.Println("updateHpa")

	var hpa apiobj.Hpa
	c.ShouldBind(&hpa)
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdHpas+"/%s/%s", namespace, name)

	hpaJson, err := json.Marshal(hpa)
	if err != nil {
		c.JSON(500, gin.H{"update": "fail"})
	}
	err = etcd.EtcdKV.Put(key, hpaJson)
	if err != nil {
		c.JSON(500, gin.H{"update": "fail"})
	}
	c.JSON(200, gin.H{"update": "success"})
}

func GetHpa(c *gin.Context) {
	fmt.Println("getHpa")
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdHpas+"/%s/%s", namespace, name)
	res, err := etcd.EtcdKV.Get(key)
	if err != nil {
		c.JSON(500, gin.H{"get": "fail"})
	}
	c.JSON(200, gin.H{"data": string(res)})
}

func GetHpaStatus(c *gin.Context) {
	fmt.Println("getHpaStatus")
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdHpas+"/%s/%s/status", namespace, name)
	res, err := etcd.EtcdKV.Get(key)
	if err != nil {
		c.JSON(500, gin.H{"get": "fail"})
	}
	c.JSON(200, gin.H{"data": string(res)})
}