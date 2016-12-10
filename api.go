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

	return getBodyFromResponse(resp)
}

func getBodyFromResponse(resp *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return body, err
}
