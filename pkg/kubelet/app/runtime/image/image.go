package image

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

/*
	PullImage
	参数：容器镜像地址
	返回：本地镜像ID，error
*/

func PullImage(imageRef string) (string, error) {
	imageIds, err := findLocalImageIdByImageRef(imageRef)
	if len(imageIds) == 1 {
		fmt.Printf("ImageCount has already been pulled")
		return imageIds[0], nil
	} else if len(imageIds) > 1 {
		for _, imageId := range imageIds {
			fmt.Println(imageId + "\n")
		}
		return "", errors.New("image count")
	}
	ctx := context.Background()
	tmpClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}
	defer tmpClient.Close()
	image, err := tmpClient.ImagePull(ctx, imageRef, image.PullOptions{})
	if err != nil {
		return "", err
	}
	defer image.Close()
	io.Copy(os.Stdout, image)
	imageIds, err = findLocalImageIdByImageRef(imageRef)
	if len(imageIds) != 1 {
		fmt.Printf("ImageCount is: %d\n", len(imageIds))
		for _, imageId := range imageIds {
			fmt.Println(imageId + "\n")
		}
		return "", errors.New("image count")
	}
	return imageIds[0], nil
}

/*
	RemoveImage
	参数：容器镜像地址
	返回：error
*/

func RemoveImage(imageRef string) error {
	ctx := context.Background()
	tmpClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer tmpClient.Close()
	imageIds, err := findLocalImageIdByImageRef(imageRef)
	if len(imageIds) != 1 {
		return errors.New("image count ")
	}
	_, err = tmpClient.ImageRemove(ctx, imageIds[0], image.RemoveOptions{})
	if err != nil {
		return err
	}
	return nil
}

func findLocalImageIdByImageRef(imageRef string) ([]string, error) {
	ctx := context.Background()
	tmpClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return []string{}, err
	}
	defer tmpClient.Close()
	filterArgs := filters.NewArgs()
	filterArgs.Add("reference", parseImageRef(imageRef))
	images, err := tmpClient.ImageList(ctx, image.ListOptions{
		Filters: filterArgs,
	})
	if err != nil {
		return []string{}, err
	}
	imageIds := []string{}
	for _, img := range images {
		imageIds = append(imageIds, img.ID)
	}
	return imageIds, nil
}

func parseImageRef(imageRef string) string {
	if !strings.HasPrefix(imageRef, "docker.io/") {
		return imageRef
	}
	parts := strings.Split(imageRef, "/")
	lastPart := parts[len(parts)-1]
	return lastPart
}
