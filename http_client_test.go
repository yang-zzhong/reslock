package reslock

import (
	"bytes"
	"errors"
	"net/http"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

type testbody struct {
	bytes.Buffer
}

func (testbody) Close() error {
	return nil
}

func testcli(t *testing.T, statusCode int, method string) (*httpClient, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	doer := NewMockrequestDoer(ctrl)
	doer.EXPECT().Do(gomock.Any()).DoAndReturn(func(req *http.Request) (res *http.Response, err error) {
		if req.URL.Path != "/127.0.0.1:8081/keys" {
			err = errors.New("path error")
			return
		}
		if req.Method != method {
			err = errors.New("method error")
			return
		}
		if req.Header.Get("Content-Type") != "application/json" {
			err = errors.New("content type error")
			return
		}
		bs := make([]byte, 512)
		var l int
		if l, err = req.Body.Read(bs); err != nil {
			return
		}
		bs = bs[:l]
		if string(bs) != "{\"keys\":[\"hello\"]}" {
			err = errors.New("body error")
			return
		}
		res = &http.Response{StatusCode: statusCode, Body: &testbody{Buffer: *bytes.NewBuffer([]byte("{\"msg\":\"error\"}"))}}
		return
	})
	cli := &httpClient{baseUrl: "http://127.0.0.1:8081", doer: doer}
	return cli, ctrl
}

func TestHttpClientLockOK(t *testing.T) {
	cli, ctrl := testcli(t, http.StatusOK, http.MethodPost)
	defer ctrl.Finish()
	if err := cli.Lock("hello"); err != nil {
		t.Fatal(err)
	}

}

func TestHttpClientLockFailed(t *testing.T) {
	cli, ctrl := testcli(t, http.StatusInternalServerError, http.MethodPost)
	defer ctrl.Finish()
	if err := cli.Lock("hello"); err == nil || err != nil && err.Error() != "error" {
		t.Fatal(err)
	}
}

func TestHttpClientUnLockOK(t *testing.T) {
	cli, ctrl := testcli(t, http.StatusOK, http.MethodDelete)
	defer ctrl.Finish()
	if err := cli.Unlock("hello"); err != nil {
		t.Fatal(err)
	}

}

func TestHttpClientUnLockFailed(t *testing.T) {
	cli, ctrl := testcli(t, http.StatusInternalServerError, http.MethodDelete)
	defer ctrl.Finish()
	if err := cli.Unlock("hello"); err == nil || err != nil && err.Error() != "error" {
		t.Fatal(err)
	}
}
