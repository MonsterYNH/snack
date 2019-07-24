package util

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func GetRequest(url string, data interface{}) ([]byte, error) {
	client := &http.Client{}

	dataBytes, _ := json.Marshal(data)
	buffer := bytes.NewBuffer(dataBytes)
	req, err := http.NewRequest("GET", url, buffer)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "APPCODE 0798ec0cf1bf420ab8adca0cdfbb61b7")
}
