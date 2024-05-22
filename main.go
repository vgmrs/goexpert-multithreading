package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func getCEP1(cep string) (map[string]string, error) {
	req, err := http.NewRequest("GET", "https://brasilapi.com.br/api/cep/v1/01153000"+cep, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var data map[string]string
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	log.Printf(" GET CEP API 1: %v", data)

	return data, nil
}

func getCEP2(cep string) (map[string]string, error) {
	req, err := http.NewRequest("GET", "http://viacep.com.br/ws/"+cep+"/json/", nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var data map[string]string
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	log.Printf(" GET CEP API 2: %v", data)

	return data, nil
}

func getCEP(cep string) (string, error) {
	go getCEP1(cep)
	getCEP2(cep)

	// TODO: use select/case to catch first response from APIs

	return cep, nil
}

func main() {
	// TODO: get cep from command line
	cep := ""

	resp, err := getCEP(cep)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}
