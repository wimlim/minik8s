package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	// NFS 挂载路径
	nfsPath := "192.168.1.14:/nfs"
	containerMountPath := "/mnt/nfs"

	config := &container.Config{
		Image: "nginx:latest", // 替换为你要使用的镜像名
	}

	hostConfig := &container.HostConfig{
		Binds: []string{
			fmt.Sprintf("%s:%s", nfsPath, containerMountPath),
		},
	}

	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, nil, "your_container_name")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}

	fmt.Printf("Container %s started\n", resp.ID)
}