package reslock

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"path"
)

type httpClient struct {
	baseUrl string
	doer    requestDoer
}

type requestDoer interface {
	Do(req *http.Request) (res *http.Response, err error)
}

type LockHttpClient interface {
	Lock(key ...string) error
	Unlock(key ...string) error
	Locked(key string) (error, bool)
}

var _ LockHttpClient = &httpClient{}

func HttpClient(baseUrl string, doer requestDoer) *httpClient {
	return &httpClient{baseUrl: baseUrl, doer: doer}
}

func (client *httpClient) Lock(key ...string) error {
	return client.lockOrUnlock(http.MethodPost, key...)
}

func (client *httpClient) Unlock(key ...string) error {
	return client.lockOrUnlock(http.MethodDelete, key...)
}

func (client *httpClient) Locked(key string) (err error, locked bool) {
	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, client.url("/keys/"+key), &bytes.Buffer{})
	if err != nil {
		return
	}
	var res *http.Response
	res, err = client.doer.Do(req)
	if err != nil {
		return
	}
	var bs []byte
	if _, err = res.Body.Read(bs); err != nil {
		return
	}
	var body struct {
		Msg    string `json:"msg"`
		Locked bool   `json:"locked"`
	}
	if err = json.Unmarshal(bs, &body); err != nil {
		return
	}
	if res.StatusCode != http.StatusOK {
		err = errors.New(body.Msg)
	}
	locked = body.Locked
	return
}

func (client *httpClient) lockOrUnlock(method string, key ...string) error {
	keys := map[string][]string{
		"keys": key,
	}
	bs, _ := json.Marshal(keys)
	buf := bytes.NewReader(bs)
	req, err := http.NewRequest(method, client.url("/keys"), buf)
	defer req.Body.Close()
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return err
	}
	res, err := client.doer.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		bs, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		var body struct {
			Msg string `json:"msg"`
		}
		if err := json.Unmarshal(bs, &body); err != nil {
			return err
		}
		return errors.New(body.Msg)
	}
	return nil
}

func (client *httpClient) url(p string) string {
	return path.Join(client.baseUrl, p)
}
