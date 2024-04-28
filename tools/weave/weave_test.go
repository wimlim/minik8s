package weave

import (
	"os/exec"
	"testing"
)

func TestWeaveAttach(t *testing.T) {
	containerID := "TestWeaveAttach"
	setupCmd := exec.Command("docker", "run", "-d", "--name", containerID, "alpine", "sleep", "3600")
	if err := setupCmd.Run(); err != nil {
		t.Fatalf("Failed to create and start test container: %s", err)
	}

	_, err := WeaveAttach(containerID)
	if err != nil {
		t.Errorf("Error attaching weave to container: %s", err)
	}

	tearDownCmd := exec.Command("docker", "rm", "-f", containerID)
	if err := tearDownCmd.Run(); err != nil {
		t.Errorf("Failed to remove test container: %s", err)
	}
}

func TestWeaveFindIpByContainerID(t *testing.T) {
	containerID := "test1"
	expectedIP := "10.32.0.1"
	ip, err := WeaveFindIpByContainerID(containerID)
	if err != nil {
		t.Errorf("Error finding IP for container: %s", err)
	}
	if ip != expectedIP {
		t.Errorf("IP did not match: got %v want %v", ip, expectedIP)
	}
}
