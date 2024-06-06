package apiconfig

import (
	"fmt"
	"testing"
)

func TestGetMasterIP (t *testing.T) {
	ip := GetMasterIP()
	fmt.Println(ip)
	fmt.Printf("%d",ServerDefaultPort)
}