package apiserverutil

import (
	"minik8s/pkg/apiobj"
	"minik8s/pkg/apirequest"
)

func GetAllNodes() ([]apiobj.Node, error) {
	return apirequest.GetAllNodes()
}
