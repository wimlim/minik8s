package handler

import(
	"fmt"
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
}

func DeletePod(c *gin.Context){
	fmt.Println("deletePod")
}

func UpdatePod(c *gin.Context){
	fmt.Println("updatePod")
}

func GetPod(c *gin.Context){
	fmt.Println("getPod")
}

func GetPodStatus(c *gin.Context){
	fmt.Println("getPodStatus")
}