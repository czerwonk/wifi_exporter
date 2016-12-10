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
	"encoding/json"
	"errors"
	"fmt"
)

type accessPoint struct {
	name     string
	state    int
	clientsN int
	clientsG int
}

type apJsonData struct {
	Data []struct {
		Name     string `json:"name"`
		State    int    `json:"state"`
		NaNumSta int    `json:"na-num_sta"`
		NgNumSta int    `json:"ng-num_sta"`
	} `json:"data"`
	Meta struct {
		Rc  string `json:"rc"`
		Msg string `json:"msg,omitempty"`
	} `json:"meta"`
}

func getAccessPoints(siteId string, cookie string) ([]*accessPoint, error) {
	ressource := fmt.Sprintf("s/%s/stat/device", siteId)
	body, err := requestApi(ressource, cookie)

	if err != nil {
		return nil, err
	}

	return parseAccessPointJson(body)
}

func parseAccessPointJson(body []byte) ([]*accessPoint, error) {
	var data apJsonData
	err := json.Unmarshal(body, &data)

	if err != nil {
		return nil, err
	}

	if data.Meta.Rc != "ok" {
		return nil, errors.New(data.Meta.Msg)
	}

	return getApsFromJsonData(data), nil
}

func getApsFromJsonData(json apJsonData) []*accessPoint {
	res := make([]*accessPoint, 0)

	for _, x := range json.Data {
		s := accessPoint{name: x.Name, state: x.State, clientsN: x.NaNumSta, clientsG: x.NgNumSta}
		res = append(res, &s)
	}

	return res
}
