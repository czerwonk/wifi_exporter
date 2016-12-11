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
