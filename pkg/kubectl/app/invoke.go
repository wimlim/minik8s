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

	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)

var invokeCmd = &cobra.Command{
	Use:   "invoke",
	Short: "Invoke a function",
	Run:   invokeHandler,
}

type Event struct {
	Type     string `yaml:"type" json:"type"`
	FilePath string `yaml:"filePath" json:"filePath"`
	FuncName string `yaml:"funcName" json:"funcName"`
	RequBody string `yaml:"requBody" json:"requBody"`
}

// invoke filepath funcname body
func invokeHandler(cmd *cobra.Command, args []string) {
	fd, err := os.Open(args[0])
	if err != nil {
		fmt.Println("open file error")
		return
	}
	defer fd.Close()
	content, err := io.ReadAll(fd)
	if err != nil {
		fmt.Println("read file error")
		return
	}

	var event Event
	err = yaml.Unmarshal(content, &event)
	if err != nil {
		fmt.Println("unmarshal error")
		return
	}
	triggerChan := make(chan bool)

	filePath := event.FilePath
	func_name := event.FuncName
	requ_body := event.RequBody

	monitorFile(filePath, func_name, requ_body, triggerChan)
	
}

func monitorFile(filePath string, func_name string, requ_body string, triggerChan <-chan bool) {
	lastModTime := getLastModifiedTime(filePath)

	for {
		select {
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
