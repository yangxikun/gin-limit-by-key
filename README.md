# Description
An in-memory Gin middleware to limit access rate by custom key.

It depends on two library:

* [golang.org/x/time/rate](https://godoc.org/golang.org/x/time/rate): rate limit
* [github.com/patrickmn/go-cache](https://github.com/patrickmn/go-cache): expire limiter related key

# installation

```
go get -u github.com/yangxikun/gin-limit-by-key
```

# Example

```go
package main

import (
    limit "github.com/yangxikun/gin-limit-by-key"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func main() {
	r := gin.Default()

	r.Use(limit.NewRateLimiter(func(c *gin.Context) string {
		return c.ClientIP() // limit rate by client ip
	}, func(c *gin.Context) (*rate.Limiter, time.Duration) {
		return rate.NewLimiter(rate.Every(time.Second), 10), time.Hour // limit 10 qps/clientIp, and the limiter liveness time duration is 1 hour
	}, func(c *gin.Context) {
		c.AbortWithStatus(429) // handle exceed rate limit request
	}))

	r.GET("/", func(c *gin.Context) {})

	r.Run(":8888")
}
```