package apirequest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)
func PostRequest(url string, apiobj interface{}) error{
	objJson, err := json.Marshal(apiobj)
	if err != nil {
		fmt.Printf("marshal  error\n")
		return nil
	}
	response, err := http.Post(url, "application/json", bytes.NewBuffer(objJson))
	if err != nil {
		fmt.Printf("post  error\n")
		return nil
	}
	defer response.Body.Close()
	return nil
}
