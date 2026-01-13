package main

import (
	"gateway/internal/middleware"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func makeProxyHandler(proxy *httputil.ReverseProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.URL.Path = c.Param("proxyPath")
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	remote, err := url.Parse(os.Getenv("AUCTION_SERVICE_URL"))
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)

	proxy.Transport = &http.Transport{
		ResponseHeaderTimeout: 500 * time.Millisecond,
	}

	r := gin.Default()
	
	r.Use(middleware.TimeoutMiddleware())

	r.Any("/*proxyPath", makeProxyHandler(proxy))

	r.Run(":8080")
}
