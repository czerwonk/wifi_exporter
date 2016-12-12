package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type loginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Strict   bool   `json:"strict"`
}

func getCookie() (string, error) {
	url := fmt.Sprintf("%s/api/login", *apiUrl)
	log.Printf("POST %s\n", url)

	d := &loginData{Username: *apiUser, Password: *apiPass, Strict: true}
	json, err := json.Marshal(d)

	if err != nil {
		log.Println(err)
		return "", err
	}

	reader := bytes.NewReader(json)
	resp, err := http.Post(url, "application/json", reader)

	if err != nil {
		log.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	return handleLoginResponse(resp)
}

func handleLoginResponse(resp *http.Response) (string, error) {
	if resp.StatusCode != 200 {
		return "", errors.New("Login failed!")
	}

	coockieHeader := resp.Header["Set-Cookie"]

	if len(coockieHeader) == 0 {
		return "", errors.New("No coockie received!")
	}

	return strings.Split(coockieHeader[0], ";")[0], nil
}
