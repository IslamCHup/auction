package proxy

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func NewReverseProxy(fullURL string, logger *slog.Logger) *httputil.ReverseProxy {

	serviceURL, err := url.Parse(fullURL)
	if err != nil {
		logger.Error("invalid SERVICE_URL", "err", err.Error())
		os.Exit(1)
	}

	serviceProxy := httputil.NewSingleHostReverseProxy(serviceURL)

	serviceProxy.Transport = &http.Transport{
		ResponseHeaderTimeout: 500 * time.Millisecond,
	}

	serviceProxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		logger.Error("upstream service error", "err", err.Error(), "path", req.URL.Path)
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusBadGateway)
		io.WriteString(rw, `{"error":"upstream service unavailable"}`)
	}

	return serviceProxy
}

// func MakeProxyHandler(proxy *httputil.ReverseProxy) gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 		// 🔥 работаем ТОЛЬКО с оригинальным URL
// 		path := strings.TrimPrefix(c.Request.URL.Path, "/api")
// 		if path == "" {
// 			path = "/"
// 		}
// 		c.Request.URL.Path = path

// 		// прокидываем auth headers
// 		if auth := c.GetHeader("Authorization"); auth != "" {
// 			c.Request.Header.Set("Authorization", auth)
// 		}
// 		if uid := c.GetHeader("X-User-Id"); uid != "" {
// 			c.Request.Header.Set("X-User-Id", uid)
// 		}

// 		proxy.ServeHTTP(c.Writer, c.Request)
// 	}
// }

func MakeProxyHandler(proxy *httputil.ReverseProxy) gin.HandlerFunc {
	return func(c *gin.Context) {

		// ❌ НИЧЕГО не делаем с Path
		// оставляем /api/auth/register как есть

		if auth := c.GetHeader("Authorization"); auth != "" {
			c.Request.Header.Set("Authorization", auth)
		}
		if uid := c.GetHeader("X-User-Id"); uid != "" {
			c.Request.Header.Set("X-User-Id", uid)
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
