package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func requestApi(ressource string, cookie string) ([]byte, error) {
	req, err := getApiRequest(ressource, cookie)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Cookie", cookie)

	return sendApiRequest(req)
}

func getApiRequest(ressource string, cookie string) (*http.Request, error) {
	url := fmt.Sprintf("%s/api/%s", *apiUrl, ressource)
	log.Printf("GET %s\n", url)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	return req, nil
}

func sendApiRequest(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	return ioutil.ReadAll(resp.Body)
}
