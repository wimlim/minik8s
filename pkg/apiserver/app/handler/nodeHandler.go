package handler

import(
	"github.com/gin-gonic/gin"
	"fmt"
)

func GetNodes(c *gin.Context){	
	fmt.Println("getNodes")
}

func AddNode(c *gin.Context){	
	fmt.Println("addNode")
}

func DeleteNode(c *gin.Context){
	fmt.Println("deleteNode")
}

func UpdateNode(c *gin.Context){
	fmt.Println("updateNode")
}

func GetNode(c *gin.Context){	
	fmt.Print("getNode")
}

func GetNodePods(c *gin.Context){
	fmt.Println("getNodePods")
}

func GetNodeStatus(c *gin.Context){	
	fmt.Println("getNodeStatus")
}
