package scheduler

import (
	"minik8s/pkg/apirequest"
	"testing"
)

func TestGetAllNodes(t *testing.T) {
	nodes, err := apirequest.GetAllNodes()
	if err != nil {
		t.Error(err)
	}
	t.Log(nodes)
}

func TestChooseNode(t *testing.T) {
	node := chooseNode()
	t.Log(node)
}
