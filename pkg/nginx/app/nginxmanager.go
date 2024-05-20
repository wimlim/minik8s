package nginxmanager

import (
	"fmt"
	"minik8s/pkg/apiobj"
	"os"
	"os/exec"
	"strings"
)

var configPath = "/root/minik8s/nginx/default.conf"

func reloadNginx() {
	cmd := exec.Command("docker", "exec", "my-nginx-container", "nginx", "-s", "reload")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error reloading Nginx: %v\nOutput: %s\n", err, output)
	} else {
		fmt.Printf("Nginx reloaded successfully:\n%s\n", output)
	}
}

func addServerBlock(hostname string, paths []apiobj.Path) {
	// open file
	file, err := os.OpenFile(configPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// write server block
	_, err = file.WriteString("server {\n")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = file.WriteString("    listen 80;\n")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = file.WriteString("    server_name " + hostname + ";\n")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = file.WriteString("    location / {\n")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, path := range paths {
		_, err = file.WriteString("        proxy_pass " + path.SubPath + ";\n")
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
	}
	reloadNginx()
}

func deleteServerBlock(hostname string) {
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
		trimLine := strings.TrimSpace(line)
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
