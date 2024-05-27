package container

import (
	"context"
	"encoding/json"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

/*
	CPU和内存使用率计算值
	https://docs.docker.com/engine/api/v1.40/#tag/Container/operation/ContainerStats
*/

func CalcContainerCPUAndMemoryUsage(containerId string) (float64, float64, error) {
	statsInfo, err := getContainerStatus(containerId)
	if err != nil {
		return 0.0, 0.0, err
	}
	CPUUsage := calcCPUUsage(statsInfo)
	MemoryUsage := calcCPUUsage(statsInfo)
	return CPUUsage, MemoryUsage, nil
}

func getContainerStatus(containerId string) (*types.StatsJSON, error) {
	ctx := context.Background()
	tmpClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer tmpClient.Close()
	containerState, err := tmpClient.ContainerStats(ctx, containerId, false)
	if err != nil {
		return nil, err
	}
	defer containerState.Body.Close()
	decoder := json.NewDecoder(containerState.Body)
	statsInfo := &types.StatsJSON{}
	err = decoder.Decode(statsInfo)
	if err != nil {
		return nil, err
	}
	return statsInfo, nil
}

func calcCPUUsage(statsInfo *types.StatsJSON) float64 {
	cpuDelta := float64(statsInfo.CPUStats.CPUUsage.TotalUsage - statsInfo.PreCPUStats.CPUUsage.TotalUsage)
	systemCpuDelta := float64(statsInfo.CPUStats.SystemUsage - statsInfo.PreCPUStats.SystemUsage)
	numberCpus := float64(statsInfo.CPUStats.OnlineCPUs)
	if numberCpus == 0 {
		numberCpus = float64(len(statsInfo.CPUStats.CPUUsage.PercpuUsage))
	}
	if systemCpuDelta > 0.0 && numberCpus > 0.0 {
		return (cpuDelta / systemCpuDelta) * numberCpus
	}
	return 0.0
}

func calcMemoryUsage(statsInfo *types.StatsJSON) float64 {
	used_memory := float64(statsInfo.MemoryStats.Usage)
	available_memory := float64(statsInfo.MemoryStats.Limit)
	if used_memory > 0.0 && available_memory > 0.0 {
		return used_memory / available_memory
	}
	return 0.0
}
