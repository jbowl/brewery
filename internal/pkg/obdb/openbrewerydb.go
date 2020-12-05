package obdb

import (
	"net/http"
	"time"
)

type OBDB struct {
	APIUrl string
}

func (s *OBDB) RESTReq(method string, url string) (*http.Response, error) {
	var client = &http.Client{
		Timeout: time.Second * 120,
	}

	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	//req.Header.Add("Authorization", "Token "+s.APIToken)

	return client.Do(req)
}
