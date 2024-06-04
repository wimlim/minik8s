package weave

import (
	"errors"
	"os/exec"
	"regexp"
	"strings"
)

// execWeaveCommand executes a weave command with the given arguments and logs the output.
func execWeaveCommand(args ...string) (string, error) {
	if len(args) == 0 || args[0] == "" {
		return "", errors.New("invalid command arguments")
	}

	out, err := exec.Command("weave", args...).CombinedOutput()
	output := string(out)
	if err != nil {
		return output, err
	}
	result := strings.TrimSuffix(output, "\r\n")
	return result, nil
}

// WeaveAttach attaches a weave network to a container.
func WeaveAttach(containerID string) (string, error) {
	if containerID == "" {
		return "", errors.New("containerID is empty")
	}
	return execWeaveCommand("attach", containerID)
}

// WeaveFindIpByContainerID finds the IP address of a container using weave.
func WeaveFindIpByContainerID(containerID string) (string, error) {
	if containerID == "" {
		return "", errors.New("containerID is empty")
	}

	output, err := execWeaveCommand("ps", containerID)
	if err != nil {
		return "", err
	}

	ipRegex := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
	matches := ipRegex.FindStringSubmatch(output)
	if len(matches) < 2 {
		return "", errors.New("could not find IP address in command output")
	}

	return matches[1], nil
}

func WeaveConnect(serverIP string) (string, error) {
	if serverIP == "" {
		return "", errors.New("serverIP is empty")
	}
	return execWeaveCommand("connect", serverIP)
}
