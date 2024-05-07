package image

import (
	"fmt"
	"testing"
)

func TestPullImage(t *testing.T) {
	imageId, err := PullImage("docker.io/library/alpine")
	fmt.Println(imageId)
	if err != nil {
		t.Error(err)
	}
}

func TestFindLocalImageIdByImageRef(t *testing.T) {
	imageIds, err := findLocalImageIdByImageRef("docker.io/library/alpine")
	fmt.Println(len(imageIds))
	if err != nil {
		t.Error(err)
	}
}
