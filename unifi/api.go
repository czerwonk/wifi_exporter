package unifi

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/prometheus/common/log"
)

func requestApi(ressource string, cookie string, url string) ([]byte, error) {
	req, err := getApiRequest(ressource, url)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Cookie", cookie)

	return sendApiRequest(req)
}

func getApiRequest(ressource string, url string) (*http.Request, error) {
	u := fmt.Sprintf("%s/api/%s", url, ressource)
	log.Debugf("GET %s\n", u)
	req, err := http.NewRequest("GET", u, nil)

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
