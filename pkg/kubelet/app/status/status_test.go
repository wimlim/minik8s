package status

import (
	"fmt"
	"minik8s/pkg/apiobj"
	"testing"
	"time"
)

// Test function for GetNodeStatus
func TestGetNodeStatus(t *testing.T) {
	// Call the function to test
	nodeStatus, err := GetNodeStatus()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check the Condition is Ready
	if nodeStatus.Condition != apiobj.NodeCondition(apiobj.Ready) {
		t.Errorf("Expected condition %v, got %v", apiobj.NodeCondition(apiobj.Ready), nodeStatus.Condition)
	}

	// Check the CpuPercent is within valid range
	if nodeStatus.CpuPercent < 0.0 || nodeStatus.CpuPercent > 100.0 {
		t.Errorf("Expected CPU percent between 0.0 and 100.0, got %v", nodeStatus.CpuPercent)
	}

	// Check the MemPercent is within valid range
	if nodeStatus.MemPercent < 0.0 || nodeStatus.MemPercent > 100.0 {
		t.Errorf("Expected Memory percent between 0.0 and 100.0, got %v", nodeStatus.MemPercent)
	}

	// Check the PodNum is non-negative
	if nodeStatus.PodNum < 0 {
		t.Errorf("Expected Pod num to be non-negative, got %v", nodeStatus.PodNum)
	}

	// Check the UpdateTime is not in the future
	if nodeStatus.UpdateTime.After(time.Now()) {
		t.Errorf("UpdateTime %v is after now", nodeStatus.UpdateTime)
	}

	// Print the result if no errors
	fmt.Printf("NodeStatus: %+v\n", *nodeStatus)
}
