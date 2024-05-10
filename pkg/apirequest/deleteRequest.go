package apirequest

import (
	"fmt"
	"net/http"
)

func DeleteRequest(url string) {
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Println("new request error")
		return
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("do request error")
		return
	}
	defer resp.Body.Close()
}
