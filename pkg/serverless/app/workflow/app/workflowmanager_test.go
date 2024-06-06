package workflow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestStartWorkflow(t *testing.T) {
	// query as curl -X POST -d '{"a":1, "b":2}' http://localhost:8081/namespaces/default/functions/example-function
	url := "http://localhost:8081/default/example-function"
	requestBody, _ := json.Marshal(map[string]int{"a": 1, "b": 2})

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Function invocation failed")
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	t.Log("hahaha")
	t.Log(string(body))
}
