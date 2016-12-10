/*
Copyright 2016 Daniel Czerwonk (d.czerwonk@gmail.com)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/

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
