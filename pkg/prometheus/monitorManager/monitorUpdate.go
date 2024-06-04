package monitormanager

import (
	"encoding/json"
	"fmt"
	"os"
)

func AddPodMonitorDataToFile(monitorData MonitorData) {
	filename := PodMonitorFilePath
	data, err := readJSONFile(filename)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}
	data = append(data, monitorData)
	err = writeJSONFile(filename, data)
	if err != nil {
		fmt.Println("Error writing JSON file:", err)
		return
	}
}

func RemovePodMonitorDataToFile(labels Labels) {
	filename := PodMonitorFilePath
	data, err := readJSONFile(filename)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}
	var filteredData []MonitorData
	for _, monitorData := range data {
		if monitorData.Labels.Name == labels.Name &&
			monitorData.Labels.Namespace == labels.Namespace &&
			monitorData.Labels.UID == labels.UID {
			continue // 跳过需要删除的项
		}
		filteredData = append(filteredData, monitorData)
	}
	err = writeJSONFile(filename, filteredData)
	if err != nil {
		fmt.Println("Error writing JSON file:", err)
		return
	}
}

func readJSONFile(filename string) ([]MonitorData, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data []MonitorData
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func writeJSONFile(filename string, data []MonitorData) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	return nil
}
