package registry

import (
	"fmt"
)

func test(){

	r := NewRegistry()
	if r == nil {
		fmt.Println("NewRegistry error")
	}

	
	r.BuildImage()

}