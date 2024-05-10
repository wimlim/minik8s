package main

import (
	"minik8s/pkg/controller/app/ctlmanager"
)
func main() {
	ctlmanager.NewControllerManager().Run(make(<-chan struct{}))
}