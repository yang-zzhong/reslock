package reslock

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

type responseWriter struct {
	header http.Header
	io     bytes.Buffer
	code   int
}

func (rw *responseWriter) Header() http.Header {
	return rw.header
}

func (rw *responseWriter) Write(bs []byte) (int, error) {
	l, err := rw.io.Write(bs)
	return int(l), err
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.code = code
}

func TestGinLocker_LockBatch(t *testing.T) {
	keys := map[string][]string{
		"keys": {"hello"},
	}
	bs, _ := json.Marshal(keys)
	buf := bytes.NewBuffer(bs)
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1/keys", buf)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	writer := &responseWriter{header: make(http.Header), io: bytes.Buffer{}}
	ctx, _ := gin.CreateTestContext(writer)
	ctx.Request = req
	ginlocker := GinLocker(New())
	ginlocker.LockBatch(ctx)
	if ginlocker.underlying.keys["hello"] != 1 {
		t.Fatal("lock hello error")
	}
	if writer.code != 200 {
		t.Fatal("gin http code error")
	}
	if writer.header.Get("Content-Type") != "application/json; charset=utf-8" {
		t.Fatal("gin header error")
	}
	if writer.io.String() != "{\"msg\":\"lock succeeded ^_^\"}" {
		t.Fatal("gin body error")
	}
}

func TestGinLocker_UnlockBatch(t *testing.T) {
	keys := map[string][]string{
		"keys": {"hello"},
	}
	bs, _ := json.Marshal(keys)
	buf := bytes.NewBuffer(bs)
	req, err := http.NewRequest(http.MethodDelete, "http://127.0.0.1/keys", buf)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	writer := &responseWriter{header: make(http.Header), io: bytes.Buffer{}}
	ctx, _ := gin.CreateTestContext(writer)
	ctx.Request = req
	locker := &Locker{
		keys: map[string]int{
			"hello": 1,
		},
	}
	ginlocker := GinLocker(locker)
	ginlocker.UnlockBatch(ctx)
	if _, ok := ginlocker.underlying.keys["hello"]; ok {
		t.Fatal("lock hello error")
	}
	if writer.code != 200 {
		t.Fatal("gin http code error")
	}
	if writer.header.Get("Content-Type") != "application/json; charset=utf-8" {
		t.Fatal("gin header error")
	}
	if writer.io.String() != "{\"msg\":\"unlock succeeded ^_^\"}" {
		t.Fatal("gin body error")
	}
}
