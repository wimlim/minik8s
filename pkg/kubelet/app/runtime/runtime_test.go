package runtime

import (
	"fmt"
	"testing"
)

func TestFindAvailablePort(t *testing.T) {
	pausePortSet := map[string]struct{}{}
	for i := 1; i <= 100; i++ {
		port, err := findAvailablePort(&pausePortSet)
		if err != nil {
			fmt.Printf("Error finding available port: %v\n", err)
			return
		}
		pausePortSet[port] = struct{}{}
		fmt.Printf("Test %d: Available Non-repeating port: %s\n", i, port)
	}
	fmt.Printf("\npausePortSet:\n")
	for port, _ := range pausePortSet {
		fmt.Printf("%s\n", port)
	}
}
