package container

import (
	"fmt"
	"minik8s/pkg/minik8sTypes"
	"testing"

	"github.com/docker/docker/api/types/filters"
)

/*
	id					Container's ID
	name				Container's name
	label				An arbitrary string representing either a key or a key-value pair. Expressed as <key> or <key>=<value>
	exited				An integer representing the container's exit code. Only useful with --all.
	status				One of created, restarting, running, removing, paused, exited, or dead
	ancestor			Filters containers which share a given image as an ancestor. Expressed as <image-name>[:<tag>], <image id>, or <image@digest>
						before or since	Filters containers created before or after a given container ID or name
	volume				Filters running containers which have mounted a given volume or bind mount.
	network				Filters running containers connected to a given network.
	publish or expose	Filters containers which publish or expose a given port. Expressed as <port>[/<proto>] or <startport-endport>/[<proto>]
	health				Filters containers based on their healthcheck status. One of starting, healthy, unhealthy or none.
	isolation			Windows daemon only. One of default, process, or hyperv.
	is-task				Filters containers that are a "task" for a service. Boolean option (true or false)
*/

const (
	IMAGE_IN_FILTER = "ancestor"
	NAME_IN_FILTER  = "name"
	ID_IN_FILTER    = "id"
)

func TestCreateContainer(t *testing.T) {
	config := minik8sTypes.ContainerConfig{
		Image: "alpine",
		Cmd:   []string{"echo", "hello world"},
		Tty:   false,
		Name:  "",
	}
	id, err := CreateContainer(config)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(id)
}

func TestListAllContainer(t *testing.T) {
	containers, err := ListAllContainer()
	if err != nil {
		t.Error(err)
	}
	for _, ctn := range containers {
		fmt.Println(ctn.ID + "\t" + ctn.Image + "\t" + ctn.Names[0])
	}
}

func TestListContainerWithOption(t *testing.T) {
	filterArgs := filters.NewArgs()
	filterArgs.Add(IMAGE_IN_FILTER, "alpine")
	containers, err := ListContainerWithFilters(filterArgs)
	if err != nil {
		t.Error(err)
	}
	for _, ctn := range containers {
		fmt.Println(ctn.ID + "\t" + ctn.Image + "\t" + ctn.Names[0])
	}
}

func TestStartContainer(t *testing.T) {
	filterArgs := filters.NewArgs()
	filterArgs.Add(IMAGE_IN_FILTER, "alpine")
	containers, err := ListContainerWithFilters(filterArgs)
	if err != nil {
		t.Error(err)
	}
	for _, ctn := range containers {
		_, err := StartContainer(ctn.ID)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestStopContainer(t *testing.T) {
	filterArgs := filters.NewArgs()
	filterArgs.Add(IMAGE_IN_FILTER, "alpine")
	containers, err := ListContainerWithFilters(filterArgs)
	if err != nil {
		t.Error(err)
	}
	for _, ctn := range containers {
		_, err := StopContainer(ctn.ID)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestRemoveContainer(t *testing.T) {
	filterArgs := filters.NewArgs()
	filterArgs.Add(IMAGE_IN_FILTER, "alpine")
	containers, err := ListContainerWithFilters(filterArgs)
	if err != nil {
		t.Error(err)
	}
	for _, ctn := range containers {
		_, err := RemoveContainer(ctn.ID)
		if err != nil {
			t.Error(err)
		}
	}
}
