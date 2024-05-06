package scheduler

import (
	"testing"
)

func TestGetAllNodes(t *testing.T) {
	nodes, err := getAllNodes()
	if err != nil {
		t.Error(err)
	}
	t.Log(nodes)
}

func TestChooseNode(t *testing.T) {
	node := chooseNode()
	t.Log(node)
}
