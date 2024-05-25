package apirequest

import (
	"bytes"

	"encoding/json"
	"fmt"
	"net/http"
)

func PutRequest(url string, apiobj interface{}) error {
	jsonData, err := json.Marshal(apiobj)
	if err != nil {
		fmt.Println("marshal  error")
		return nil
	}
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("create put request error:", err)
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("put error:", err)
		return nil
	}
	defer response.Body.Close()
	return nil
}
