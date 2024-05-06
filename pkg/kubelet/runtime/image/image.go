package image

import (
	"context"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"io"
	"os"
)

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
	return imageRef, nil
}
