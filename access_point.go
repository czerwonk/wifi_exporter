package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ssid struct {
	name     string
	clientsN int
	clientsG int
}

type accessPoint struct {
	name     string
	mac      string
	state    int
	clientsN int
	clientsG int
	ssids    []*ssid
}

type apJsonData struct {
	Name     string `json:"name,omitempty"`
	State    int    `json:"state"`
	NaNumSta int    `json:"na-num_sta"`
	NgNumSta int    `json:"ng-num_sta"`
	Mac      string `json:"mac"`
	VapTable []struct {
		Essid  string `json:"essid"`
		Radio  string `json:"radio"`
		NumSta int    `json:"num_sta"`
	} `json:"vap_table"`
}

type apJson struct {
	Data []apJsonData `json:"data"`
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
	var data apJson
	err := json.Unmarshal(body, &data)

	if err != nil {
		return nil, err
	}

	if data.Meta.Rc != "ok" {
		return nil, errors.New(data.Meta.Msg)
	}

	return getApsFromJsonData(data), nil
}

func getApsFromJsonData(json apJson) []*accessPoint {
	res := make([]*accessPoint, 0)

	for _, x := range json.Data {
		s := accessPoint{name: x.Name, mac: x.Mac, state: x.State, clientsN: x.NaNumSta, clientsG: x.NgNumSta, ssids: getSsids(x)}
		res = append(res, &s)
	}

	return res
}

func getSsids(data apJsonData) []*ssid {
	m := make(map[string]*ssid)

	for _, x := range data.VapTable {
		s := getOrAddSsid(x.Essid, m)

		if x.Radio == "ng" {
			s.clientsG += x.NumSta
		} else {
			s.clientsN += x.NumSta
		}
	}

	ssids := make([]*ssid, 0)

	for _, v := range m {
		ssids = append(ssids, v)
	}

	return ssids
}

func getOrAddSsid(name string, ssids map[string]*ssid) *ssid {
	if s, ok := ssids[name]; ok {
		return s
	}

	s := &ssid{name: name}
	ssids[name] = s

	return s
}
