package host

import (
	"fmt"
	"testing"
)

func TestGetHostCPUPercent(t *testing.T) {
	cpuPercent, err := GetHostCPUPercent()
	if err != nil {
		t.Errorf("GetHostCPUPercent() failed with error: %v", err)
	}

	if cpuPercent < 0 || cpuPercent > 100 {
		t.Errorf("CPU percent out of range: %.2f", cpuPercent)
	} else {
		fmt.Printf("CPU percent: %.2f\n", cpuPercent)
	}
}

func TestGetHostMemoryPercent(t *testing.T) {
	memPercent, err := GetHostMemoryPercent()
	if err != nil {
		t.Errorf("GetHostMemoryPercent() failed with error: %v", err)
	}

	if memPercent < 0 || memPercent > 100 {
		t.Errorf("Memory percent out of range: %.2f", memPercent)
	} else {
		fmt.Printf("Memory percent: %.2f\n", memPercent)
	}
}

func TestGetHostIP(t *testing.T) {
	ip, err := GetHostIP()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Host IP:", ip)
	}
}
