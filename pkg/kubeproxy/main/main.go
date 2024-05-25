package main

import kubeproxy "minik8s/pkg/kubeproxy/app"

func main() {
	kp := kubeproxy.NewKubeProxy()
	kp.Run()
}
