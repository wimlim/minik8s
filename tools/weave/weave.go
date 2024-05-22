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

// WeaveExpose exposes a service using weave.
func WeaveExpose(serviceID string) error {
	if serviceID == "" {
		return errors.New("serviceID is empty")
	}
	_, err := execWeaveCommand("expose", serviceID+"/32")
	return err
}

// WeaveHide hides a service using weave.
func WeaveHide(serviceID string) error {
	if serviceID == "" {
		return errors.New("serviceID is empty")
	}
	_, err := execWeaveCommand("hide", serviceID)
	return err
}

// WeaveDNSAdd adds a DNS entry using weave.
func WeaveDNSAdd(hostname string, ip string) error {
	if hostname == "" {
		return errors.New("hostname is empty")
	}
	if ip == "" {
		return errors.New("ip is empty")
	}
	_, err := execWeaveCommand("dns-add", hostname, ip)
	return err
}

// WeaveDNSRemove removes a DNS entry using weave.
func WeaveDNSRemove(hostname string) error {
	if hostname == "" {
		return errors.New("hostname is empty")
	}
	_, err := execWeaveCommand("dns-remove", hostname)
	return err
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
