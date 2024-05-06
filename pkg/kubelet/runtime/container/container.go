package container

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"minik8s/pkg/minik8sTypes"
)

func CreateContainer(containerConfig minik8sTypes.ContainerConfig) (string, error) {
	ctx := context.Background()
	tmpClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}
	defer tmpClient.Close()

	result, err := tmpClient.ContainerCreate(ctx,
		&container.Config{
			Image: containerConfig.Image,
			Cmd:   containerConfig.Cmd,
			Tty:   containerConfig.Tty,
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
