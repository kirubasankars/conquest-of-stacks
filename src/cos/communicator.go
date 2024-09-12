package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type Communicator struct {
	host   string
	apiKey string
}

func (communicator Communicator) Sent(msg string) {
	fmt.Println(msg)
	req, _ := http.NewRequest("POST", communicator.host+"/api/publish", bytes.NewBuffer([]byte(msg)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", communicator.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	fmt.Println("Response Status:", resp.Status)
}

func NewCommunicator() *Communicator {
	ws := new(Communicator)
	ws.host = "http://localhost:8000"
	//communicator.apiKey = os.Getenv("CENTRIFUGO_API_KEY")
	ws.apiKey = "9a9d85c6-9a9c-4675-85d7-27e35884a660"
	return ws
}
