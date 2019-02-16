package limit

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"runtime"
	"testing"
	"time"
)

func TestLimitByKey(t *testing.T) {
	r := gin.Default()

	r.Use(NewRateLimiter(func(c *gin.Context) string {
		return c.ClientIP()
	}, func(c *gin.Context) (*rate.Limiter, time.Duration) {
		return rate.NewLimiter(rate.Every(time.Second), 10), time.Hour
	}, func(c *gin.Context) {
		c.AbortWithStatus(429)
	}))

	r.GET("/", func(c *gin.Context) {})

	go func() {
		err := r.Run(":8888")
		if err != nil {
			t.Error("Error run http server", err.Error())
		}
	}()

	runtime.Gosched()

	for i := 0; i < 12; i++ {
		resp, err := http.DefaultClient.Get("http://127.0.0.1:8888")
		if err != nil {
			t.Error("Error during requests", err.Error())
			return
		}
		switch {
		case i < 10:
			if resp.StatusCode != 200 {
				t.Error("Unexpected status code", resp.StatusCode)
			}
		case i == 10:
			if resp.StatusCode != 429 {
				t.Error("Threashold break not detected")
			} else {
				time.Sleep(time.Second)
			}
		case i == 11:
			if resp.StatusCode == 429 {
				t.Error("Unnecessary block")
			}
		}
	}
}
