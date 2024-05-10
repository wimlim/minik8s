package container

import (
	"context"
	"minik8s/pkg/minik8sTypes"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

/*
	CreateContainer
	参数：ContainerConfig
	返回：ContainerId，error
*/

func CreateContainer(containerConfig *minik8sTypes.ContainerConfig) (string, error) {
	ctx := context.Background()
	tmpClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}
	defer tmpClient.Close()

	result, err := tmpClient.ContainerCreate(ctx,
		&container.Config{
			Image:        containerConfig.Image,
			Cmd:          containerConfig.Cmd,
			Env:          containerConfig.Env,
			Tty:          containerConfig.Tty,
			Labels:       containerConfig.Labels,
			Entrypoint:   containerConfig.Entrypoint,
			Volumes:      containerConfig.Volumes,
			ExposedPorts: containerConfig.ExposedPorts,
		},
		&container.HostConfig{
			NetworkMode:  container.NetworkMode(containerConfig.NetworkMode),
			Binds:        containerConfig.Binds,
			PortBindings: containerConfig.PortBindings,
			IpcMode:      container.IpcMode(containerConfig.IpcMode),
			PidMode:      container.PidMode(containerConfig.PidMode),
			VolumesFrom:  containerConfig.VolumesFrom,
			Links:        containerConfig.Links,
			Resources: container.Resources{
				Memory:   containerConfig.Memory,
				NanoCPUs: containerConfig.NanoCPUs,
			},
		},
		nil, nil,
		containerConfig.Name,
	)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}

/*
	RemoveContainer
	参数：ContainerId
	返回：ContainerId，error
*/

func RemoveContainer(containerId string) (string, error) {
	ctx := context.Background()
	tmpClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}
	defer tmpClient.Close()
	//先停止对应容器
	err = tmpClient.ContainerStop(ctx, containerId, container.StopOptions{})
	if err != nil {
		return "", err
	}
	//再删除容器
	err = tmpClient.ContainerRemove(ctx, containerId, container.RemoveOptions{})
	if err != nil {
		return "", err
	}
	return containerId, nil
}

/*
	StartContainer
	参数：ContainerId
	返回：ContainerId，error
*/

func StartContainer(containerId string) (string, error) {
	ctx := context.Background()
	tmpClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}
	defer tmpClient.Close()
	err = tmpClient.ContainerStart(ctx, containerId, container.StartOptions{})
	if err != nil {
		return "", err
	}
	return containerId, nil
}

/*
	StopContainer
	参数：ContainerId
	返回：ContainerId，error
*/

func StopContainer(containerId string) (string, error) {
	ctx := context.Background()
	tmpClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}
	defer tmpClient.Close()
	err = tmpClient.ContainerStop(ctx, containerId, container.StopOptions{})
	if err != nil {
		return "", err
	}

	return containerId, nil
}

func RestartContainer(containerId string) (string, error) {
	ctx := context.Background()
	tmpClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}
	defer tmpClient.Close()
	err = tmpClient.ContainerRestart(ctx, containerId, container.StopOptions{})
	if err != nil {
		return "", err
	}

	return containerId, nil
}

func ListContainerWithFilters(filterArgs filters.Args) ([]types.Container, error) {
	ctx := context.Background()
	tmpClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return []types.Container{}, err
	}
	defer tmpClient.Close()
	contianers, err := tmpClient.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: filterArgs,
	})
	if err != nil {
		return []types.Container{}, err
	}
	return contianers, nil
}

func ListAllContainer() ([]types.Container, error) {
	ctx := context.Background()
	tmpClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return []types.Container{}, err
	}
	defer tmpClient.Close()
	contianers, err := tmpClient.ContainerList(ctx, container.ListOptions{
		All: true,
	})
	if err != nil {
		return []types.Container{}, err
	}
	return contianers, nil
}

func InspectContainer(containerId string) (*types.ContainerJSON, error) {
	ctx := context.Background()
	tmpClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer tmpClient.Close()
	ctnJSON, err := tmpClient.ContainerInspect(ctx, containerId)
	if err != nil {
		return nil, err
	}
	return &ctnJSON, nil
}
