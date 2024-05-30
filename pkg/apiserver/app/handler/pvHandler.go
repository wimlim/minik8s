package handler

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/etcd"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetGlobalPVs(c *gin.Context) {
	fmt.Println("getGlobalPVs")
	key := etcd.PATH_EtcdPVs
	var resList []string
	resList, err := etcd.EtcdKV.GetPrefix(key)
	if err != nil {
		c.JSON(500, gin.H{"get": "fail"})
	}
	c.JSON(200, gin.H{
		"data": resList,
	})

}

func GetAllPVs(c *gin.Context) {
	fmt.Println("getAllPVs")
	namespace := c.Param("namespace")
	key := fmt.Sprintf(etcd.PATH_EtcdPVs+"/%s", namespace)
	var resList []string
	resList, err := etcd.EtcdKV.GetPrefix(key)
	if err != nil {
		c.JSON(500, gin.H{"get": "fail"})
	}
	c.JSON(200, gin.H{
		"data": resList,
	})
}

func AddPV(c *gin.Context) {
	fmt.Println("addPV")
	var pv apiobj.PV
	c.ShouldBind(&pv)
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdPVs+"/%s/%s", namespace, name)

	pv.MetaData.UID = uuid.New().String()[:16]

	pvJson, err := json.Marshal(pv)
	if err != nil {
		c.JSON(500, gin.H{"add": "fail"})
	}

	etcd.EtcdKV.Put(key, pvJson)
	c.JSON(200, gin.H{"add": string(pvJson)})

	mntPath := apiobj.NfsMntPath
	newPath := fmt.Sprintf("%s/%s", mntPath, pv.MetaData.Name)

	err = os.Mkdir(newPath, 0755)
	if err != nil {
		fmt.Println(err)
	}
	
}

func DeletePV(c *gin.Context) {
	fmt.Println("deletePV")
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdPVs+"/%s/%s", namespace, name)

	res, err := etcd.EtcdKV.Get(key)
	if err != nil {
		c.JSON(500, gin.H{"get": "fail"})
	}
	var pv apiobj.PV
	json.Unmarshal(res, &pv)

	mntPath := apiobj.NfsMntPath
	newPath := fmt.Sprintf("%s/%s", mntPath, pv.MetaData.Name)

	err = os.RemoveAll(newPath)
	if err != nil {
		fmt.Println(err)
	}

	etcd.EtcdKV.Delete(key)
	c.JSON(200, gin.H{"delete": "success"})
}

func UpdatePV(c *gin.Context) {
	fmt.Println("updatePV")

	var pv apiobj.PV
	c.ShouldBind(&pv)
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdPVs+"/%s/%s", namespace, name)

	pvJson, err := json.Marshal(pv)
	if err != nil {
		c.JSON(500, gin.H{"update": "fail"})
	}
	etcd.EtcdKV.Put(key, pvJson)
	c.JSON(200, gin.H{"update": "success"})
}

func GetPV(c *gin.Context) {
	fmt.Println("getPV")
	namespace := c.Param("namespace")
	name := c.Param("name")
	key := fmt.Sprintf(etcd.PATH_EtcdPVs+"/%s/%s", namespace, name)
	res, err := etcd.EtcdKV.Get(key)
	if err != nil {
		c.JSON(500, gin.H{"get": "fail"})
	}
	c.JSON(200, gin.H{
		"data": string(res),
	})
}

