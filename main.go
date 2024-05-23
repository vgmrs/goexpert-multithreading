package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type response struct {
	data    map[string]interface{}
	apiName string
}

func makeRequest(apiURL string, apiName string, ch chan<- response, errChan chan<- error) {
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		errChan <- err
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		errChan <- err
		return
	}
	defer res.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		errChan <- err
		return
	}

	ch <- response{data: data, apiName: apiName}
}

func getCEP(cep string) (response, error) {
	chan1 := make(chan response)
	chan2 := make(chan response)
	errChan := make(chan error, 2)

	go makeRequest("https://brasilapi.com.br/api/cep/v1/"+cep, "BrasilAPI", chan1, errChan)
	go makeRequest("http://viacep.com.br/ws/"+cep+"/json/", "ViaCEP", chan2, errChan)

	select {
	case resp := <-chan1:
		return resp, nil
	case resp := <-chan2:
		return resp, nil
	case err := <-errChan:
		return response{}, err
	case <-time.After(time.Second * 1):
		return response{}, fmt.Errorf("timeout")
	}
}

func main() {
	var cep string
	fmt.Print("Enter CEP number: ")
	fmt.Scanln(&cep)

	resp, err := getCEP(cep)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("Response from API %s with result: %v", resp.apiName, resp.data)
}
