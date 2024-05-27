package registry

import (
	"fmt"
	"minik8s/pkg/apiobj"
	"testing"
)

func TestBuildImage(t *testing.T){

	r := NewRegistry()
	if r == nil {
		fmt.Println("NewRegistry error")
	}

	function := apiobj.Function{
		ApiVersion: "v1",
		Kind:       "Function",
		MetaData: apiobj.MetaData{
			UID: "example-function",
		},
		Spec: apiobj.FunctionSpec{
			Path: "func.py",
			Content: []byte(`def main(params):
	a = params["a"]
	b = params["b"]
	c = a * b
	return c`),
		},
	}
	
	r.BuildImage(function)

}