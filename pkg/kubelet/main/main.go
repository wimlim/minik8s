package main

import (
	kubelet "minik8s/pkg/kubelet/app"
)

func main() {
	kubelet.Run()
}