package host

import (
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func GetHostIP() (string, error) {
	// 获取主机名
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	// 获取主机IP地址
	addrs, err := net.LookupHost(hostname)
	if err != nil {
		return "", err
	}

	// 返回第一个IP地址
	return addrs[0], nil
}

// GetHostCPUPercent 获取主机的CPU使用率（百分比）
func GetHostCPUPercent() (float64, error) {
	// 执行 ps 命令获取 CPU 使用率
	cmd := exec.Command("ps", "-A", "-o", "%cpu")
	output, err := cmd.Output()
	if err != nil {
		return 0.0, err
	}

	// 解析命令输出并计算 CPU 使用率总和
	lines := strings.Split(string(output), "\n")
	var totalCPU float64
	for _, line := range lines[1:] { // 跳过第一行标题
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		cpuPercent, err := strconv.ParseFloat(line, 64)
		if err != nil {
			return 0.0, err
		}
		totalCPU += cpuPercent
	}

	return totalCPU, nil
}

// GetHostMemoryPercent 获取主机的内存使用率（百分比）
func GetHostMemoryPercent() (float64, error) {
	// 执行 ps 命令获取内存使用率
	cmd := exec.Command("ps", "-A", "-o", "%mem")
	output, err := cmd.Output()
	if err != nil {
		return 0.0, err
	}

	// 解析命令输出并计算内存使用率总和
	lines := strings.Split(string(output), "\n")
	var totalMem float64
	for _, line := range lines[1:] { // 跳过第一行标题
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		memPercent, err := strconv.ParseFloat(line, 64)
		if err != nil {
			return 0.0, err
		}
		totalMem += memPercent
	}

	return totalMem, nil
}
