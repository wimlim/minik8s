package kubectl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"minik8s/pkg/config/serverlessconfig"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var invokeCmd = &cobra.Command{
	Use:   "invoke",
	Short: "Invoke a function",
	Run:   invokeHandler,
}

// invoke filepath funcname body
func invokeHandler(cmd *cobra.Command, args []string) {
	filePath := args[0]  // 替换为您要监控的文件路径
	func_name := args[1] // 替换为您要发送的请求体
	requ_body := args[2] // 替换为您要发送的请求体
	stopChan := make(chan bool)
	triggerChan := make(chan bool)

	go monitorFile(filePath, func_name, requ_body, stopChan, triggerChan)

	time.Sleep(60 * time.Second)
	stopChan <- true
}

func monitorFile(filePath string, func_name string, requ_body string, stopChan, triggerChan <-chan bool) {
	lastModTime := getLastModifiedTime(filePath)

	for {
		select {
		case <-stopChan:
			fmt.Println("Stopping file monitoring.")
			return
		case <-triggerChan:
			fmt.Println("Triggering request manually.")
			sendRequest(func_name, requ_body)
		default:
			newModTime := getLastModifiedTime(filePath)
			if newModTime != lastModTime {
				fmt.Println("File has been modified!")
				sendRequest(func_name, requ_body)
				lastModTime = newModTime
			}
			time.Sleep(5 * time.Second)
		}
	}
}

func getLastModifiedTime(filePath string) time.Time {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		panic(err)
	}
	return fileInfo.ModTime()
}

func sendRequest(func_name string, requ_body string) {

	URL := serverlessconfig.GetServerlessServerUrl()
	URL = URL + "/default/" + func_name

	parts := strings.Split(requ_body, " ")
	result := make(map[string]int)
	for _, part := range parts {
		kv := strings.Split(part, ":")
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			value, _ := strconv.Atoi(strings.TrimSpace(kv[1]))
			result[key] = value
		}
	}

	requ_body_byte, _ := json.Marshal(result)
	response, err := http.Post(URL, "application/json", bytes.NewBuffer(requ_body_byte))
	if err != nil {
		fmt.Printf("post  error\n")
		return
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	fmt.Println("Response:", string(body))
}
