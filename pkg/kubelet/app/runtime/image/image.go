package image

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

/*
 * PullImage
 * 参数：容器镜像地址
 * 返回：本地镜像ID
 */

func PullImage(imageRef string) (string, error) {
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
	imageIds, err := findLocalImageIdByImageRef(imageRef)
	if len(imageIds) != 1 {
		return "", errors.New("image count ")
	}
	return imageIds[0], nil
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
	parts := strings.Split(imageRef, "/")
	lastPart := parts[len(parts)-1]
	return lastPart
}
