package reslock

import (
	"fmt"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

type ginLocker struct {
	underlying *Locker
}

func GinLocker(locker *Locker) *ginLocker {
	return &ginLocker{underlying: locker}
}

func (locker *ginLocker) Bind(engine *gin.Engine, prefix string) {
	engine.GET(path.Join(prefix, "/keys/:key"), locker.Locked)
	engine.POST(path.Join(prefix, "/keys"), locker.LockBatch)
	engine.DELETE(path.Join(prefix, "/keys"), locker.UnlockBatch)
}

func (locker *ginLocker) Lock(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "can't find key",
		})
		return
	}
	locker.underlying.Lock(key)
	c.JSON(http.StatusOK, gin.H{
		"msg": "lock succeeded ^_^",
	})
}

func (locker *ginLocker) LockBatch(c *gin.Context) {
	var form map[string][]string
	c.BindJSON(&form)
	keys, ok := form["keys"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "can't find keys",
		})
	}
	for _, key := range keys {
		locker.underlying.Lock(key)
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "lock succeeded ^_^",
	})
}

func (locker *ginLocker) UnlockBatch(c *gin.Context) {
	var form map[string][]string
	c.BindJSON(&form)
	keys, ok := form["keys"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "can't find keys",
		})
	}
	for _, key := range keys {
		locker.underlying.Unlock(key)
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "unlock succeeded ^_^",
	})
}

func (locker *ginLocker) Locked(c *gin.Context) {
	if locker.underlying.Locked(c.Param("key")) {
		c.JSON(http.StatusOK, gin.H{
			"msg":    fmt.Sprintf("key [%s] has been locked", c.Param("key")),
			"locked": true,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":    fmt.Sprintf("key [%s] has not been locked", c.Param("key")),
		"locked": false,
	})
}
