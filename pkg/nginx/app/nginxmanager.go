package nginxmanager

import (
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/apirequest"
	"os"
	"os/exec"
	"strings"
)

var configPath = "/root/minik8s/pkg/nginx/default.conf"

func reloadNginx() {
	cmd := exec.Command("docker", "exec", "my-nginx-container", "nginx", "-s", "reload")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error reloading Nginx: %v\nOutput: %s\n", err, output)
	} else {
		fmt.Printf("Nginx reloaded successfully:\n%s\n", output)
	}
}

func AddServerBlock(hostname string, paths []apiobj.Path) {
	// open file
	file, err := os.OpenFile(configPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// write server block
	file.WriteString("\nserver {\n")
	file.WriteString("    listen 80;\n")
	_, err = file.WriteString("    server_name " + hostname + ";\n")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, path := range paths {
		_, err = file.WriteString("    location " + path.SubPath + " {\n")
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = file.WriteString("        proxy_pass http://" + path.ServiceIp + ":" + fmt.Sprintf("%d", path.ServicePort) + "/;\n")
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = file.WriteString("        proxy_set_header Host $host;\n")
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = file.WriteString("        proxy_set_header X-Real-IP $remote_addr;\n")
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = file.WriteString("        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;\n")
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = file.WriteString("        proxy_set_header X-Forwarded-Proto $scheme;\n")
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = file.WriteString("    }\n")
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	_, err = file.WriteString("}\n")
	if err != nil {
		fmt.Println(err)
		return
	}
	reloadNginx()
}

func DeleteServerBlock(hostname string) {
	content, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	lines := strings.Split(string(content), "\n")
	var result []string

	inServerBlock := false
	serverBlock := []string{}

	for _, line := range lines {
		trimLine := line
		if trimLine == "server {" {
			inServerBlock = true
			serverBlock = []string{line}
			continue
		}

		if inServerBlock {
			serverBlock = append(serverBlock, line)
			if trimLine == "}" {
				inServerBlock = false
				if !containsServerName(serverBlock, hostname) {
					result = append(result, serverBlock...)
				}
				serverBlock = nil
			}
		} else {
			result = append(result, line)
		}
	}

	err = os.WriteFile(configPath, []byte(strings.Join(result, "\n")), 0644)
	if err != nil {
		fmt.Println(err)
	}
	reloadNginx()
}

func containsServerName(block []string, serverName string) bool {
	for _, line := range block {
		if strings.Contains(line, "server_name") && strings.Contains(line, serverName) {
			return true
		}
	}
	return false
}

func AddServiceIPVS(serviceSpec apiobj.ServiceSpec) {
	pods, err := apirequest.GetAllPods()
	if err != nil {
		fmt.Println("Failed to get all pods:", err)
		return
	}
	for _, port := range serviceSpec.Ports {
		var podIPs []string
		for _, pod := range pods {
			if podMatchesService(&pod, &serviceSpec) {
				podIPs = append(podIPs, pod.Status.PodIP)
			}
		}
		if len(podIPs) == 0 {
			fmt.Println("No pods match service selector")
			return
		}
		// docker exec my-nginx-container ipvsadm -A -t
		cmd := exec.Command("docker", "exec", "my-nginx-container", "ipvsadm", "-A", "-t", serviceSpec.ClusterIP+":"+fmt.Sprint(port.Port), "-s", "rr")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error adding IPVS service:", err, output)
		}
		for _, podIP := range podIPs {
			// docker exec my-nginx-container ipvsadm -a -t
			cmd := exec.Command("docker", "exec", "my-nginx-container", "ipvsadm", "-a", "-t", serviceSpec.ClusterIP+":"+fmt.Sprint(port.Port), "-r", podIP+":"+fmt.Sprint(port.Port), "-m")
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Error adding IPVS destination for pod %s on port %d: %v\nOutput: %s\n", podIP, port.Port, err, output)
			}
		}
	}
}

func podMatchesService(pod *apiobj.Pod, serviceSpec *apiobj.ServiceSpec) bool {
	labels := pod.MetaData.Labels
	for key, value := range serviceSpec.Selector {
		if currentValue, ok := labels[key]; !ok || currentValue != value {
			return false
		}
	}
	return true
}

func DeleteServiceIPVS(serviceSpec apiobj.ServiceSpec) {
	// docker exec my-nginx-container ipvsadm -D -t <clusterIP>
	cmd := exec.Command("docker", "exec", "my-nginx-container", "ipvsadm", "-D", "-t", serviceSpec.ClusterIP)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error deleting IPVS service:", err, output)
	}
}
