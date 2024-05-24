package registry

import (
	"context"
	"fmt"
	"io"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/config/serverlessconfig"
	cm "minik8s/pkg/kubelet/app/runtime/container"
	im "minik8s/pkg/kubelet/app/runtime/image"
	"minik8s/pkg/minik8sTypes"
	"os"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"

	"github.com/docker/go-connections/nat"
)

type Registry struct {
}

func NewRegistry() *Registry {
	res, err := cm.ListAllContainer()
	if err != nil {
		fmt.Println("ListAllContainer error:", err)
	}

	for _, container := range res {
		if container.Names[0] == serverlessconfig.RegistryContainerName {
			return &Registry{}
		}
	}

	im.PullImage(serverlessconfig.RegistryImage)

	port_protocol := fmt.Sprintf("%d/tcp", serverlessconfig.RegistryServerPort)
	cid, _ := cm.CreateContainer(&minik8sTypes.ContainerConfig{
		Image: serverlessconfig.RegistryImage,
		Name:  serverlessconfig.RegistryContainerName,
		PortBindings: map[nat.Port][]nat.PortBinding{nat.Port(port_protocol): {
			{HostIP: serverlessconfig.RegistryServerBindIp,
				HostPort: fmt.Sprintf("%d", serverlessconfig.RegistryServerPort),
			},
		}},

		ExposedPorts: map[nat.Port]struct{}{nat.Port(port_protocol): {}},
		Env:          []string{"REGISTRY_STORAGE_DELETE_ENABLED=true"},
	})

	cm.StartContainer(cid)
	return &Registry{}
}


func (r *Registry) PullImage(imageRef string) {
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	info, _ := cli.ImagePull(context.Background(), imageRef, image.PullOptions{
		RegistryAuth: "auth",
		All:          false,
	})

	defer info.Close()
	io.Copy(os.Stdout, info)
}

func (r *Registry) PushImage(imageRef string) {
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	info, _ := cli.ImagePush(context.Background(), imageRef, image.PushOptions{
		RegistryAuth: "auth",
		All:          false,
	})

	defer info.Close()
	io.Copy(os.Stdout, info)
}

func (r *Registry) BuildImage(f apiobj.Function) {
	os.Mkdir(f.MetaData.UID, 0777)

	curDir, _ := os.Getwd()

	

}