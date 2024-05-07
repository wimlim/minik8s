package image

import (
	"context"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"io"
	"os"
)

/*
 * PullImage
 * 参数：容器镜像地址
 * 
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
	imageIds, err := tmpClient.
	return imageRef, nil
}

func findLocalImageIdByImageRef(imageRef string) ([]string, error) {
	ctx := context.Background
	tmpClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return []string{}, err
	}
	defer tmpClient.Close()
	imageList, err := tmpClient.imageList(ctx, image.ListOptions{})
	if err != nil {
		return []string{}, err
	}
	
}