package main

import (
	"encoding/json"
	"errors"
)

type site struct {
	name string
	id   string
}

type siteJsonData struct {
	Data []struct {
		Desc string `json:"desc"`
		Name string `json:"name"`
	} `json:"data"`
	Meta struct {
		Rc  string `json:"rc"`
		Msg string `json:"msg,omitempty"`
	} `json:"meta"`
}

func getSites(cookie string) ([]*site, error) {
	body, err := requestApi("self/sites", cookie)

	if err != nil {
		return nil, err
	}

	return parseSiteJson(body)
}

func parseSiteJson(body []byte) ([]*site, error) {
	var data siteJsonData
	err := json.Unmarshal(body, &data)

	if err != nil {
		return nil, err
	}

	if data.Meta.Rc != "ok" {
		return nil, errors.New(data.Meta.Msg)
	}

	return getSitesFromJsonData(data), nil
}

func getSitesFromJsonData(json siteJsonData) []*site {
	res := make([]*site, 0)

	for _, x := range json.Data {
		s := site{name: x.Desc, id: x.Name}
		res = append(res, &s)
	}

	return res
}
