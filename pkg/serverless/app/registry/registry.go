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
	"os/exec"

	"path/filepath"

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
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println("Failed to create Docker client:", err)
		return
	}

	info, err := cli.ImagePush(context.Background(), imageRef, image.PushOptions{
		RegistryAuth: "auth",
		All:          false,
	})
	if err != nil {
		fmt.Println("Failed to push image:", err)
		return
	}

	defer info.Close()
	io.Copy(os.Stdout, info)
}

func (r *Registry) BuildImage(f apiobj.Function) {

	curpath, err := os.Getwd()
	if err != nil {
		fmt.Println("Getwd error:", err)
	}
	fmt.Println("curpath:", curpath)
	err = os.Mkdir(f.MetaData.UID, 0777)
	if err != nil {
		fmt.Println("Mkdir error:", err)
	}
	fd, err := os.Create(filepath.Join(curpath, f.MetaData.UID, "func.py"))
	if err != nil {
		fmt.Println("Create error:", err)
	}
	defer fd.Close()
	_, err = fd.Write(f.Spec.Content)
	if err != nil {
		fmt.Println("Write error:", err)
	}

	dockerfilePath := filepath.Join(curpath, f.MetaData.UID, "Dockerfile")
	Dockerfile, err := os.Create(dockerfilePath)
	if err != nil {
		fmt.Println("Create dockerfile error:", err)
	}
	Dockerfile.WriteString("FROM 10.119.13.134:5000/server_base:latest\n")
	Dockerfile.WriteString("COPY func.py /app/\n")
	Dockerfile.Close()

	imageName := fmt.Sprintf("func/%s:latest", f.MetaData.Name)
	imageRef := fmt.Sprintf("%s/%s", serverlessconfig.GetRegistryServerUrl(), imageName)

	cmd := exec.Command("docker", "build", "-t", imageRef, "-f", dockerfilePath, filepath.Join(curpath, f.MetaData.UID))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("docker build error:", err)
	}

	r.PushImage(imageRef)

	err = os.RemoveAll(filepath.Join(curpath, f.MetaData.UID))
	if err != nil {
		fmt.Println("RemoveAll error:", err)
	}
}
