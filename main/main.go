package main

import (
	"flag"
	"fmt"
	"math/rand"
	"reslock"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	var port string
	var t bool
	var prefix string
	flag.StringVar(&port, "port", "8081", "port for listening, default is 8081")
	flag.StringVar(&prefix, "prefix", "", "prefix")
	flag.BoolVar(&t, "test", false, "only test")
	r := gin.Default()
	lock := reslock.New()
	quit := make(chan bool)
	if t {
		for i := 0; i <= 10000; i++ {
			lock.Lock(fmt.Sprintf("xxx-xxx%d", i))
		}
		go func() {
			for {
				select {
				case <-quit:
					return
				default:
					start := rand.Intn(9900)
					for j := start; j < start+100; j++ {
						key := fmt.Sprintf("xxx-xxx%d", j)
						if lock.Locked(key) {
							lock.Unlock(key)
						} else {
							lock.Lock(key)
						}
					}
					time.Sleep(time.Second)
				}
			}
		}()
	}
	ginlocker := reslock.GinLocker(lock)
	ginlocker.Bind(r, prefix)
	if err := r.Run(":" + port); err != nil {
		panic(err)
	}
	quit <- true
}
