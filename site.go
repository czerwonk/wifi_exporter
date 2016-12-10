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
