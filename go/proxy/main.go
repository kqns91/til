package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	routes := map[string]gin.HandlerFunc{
		"/login":    endpoint1Handler,
		"/callback": endpoint2Handler,
	}

	gateway := r.Group("/gateway")

	for route, handler := range routes {
		gateway.Any(route, handler)
	}

	r.Any("/api/*any", reverseProxyHandler)

	r.Run(":8082")
}

func endpoint1Handler(c *gin.Context) {
	c.String(http.StatusOK, "This is endpoint 1 under gateway.")
}

func endpoint2Handler(c *gin.Context) {
	c.String(http.StatusOK, "This is endpoint 2 under gateway.")
}

func reverseProxyHandler(c *gin.Context) {
	var targetURL *url.URL

	switch c.Request.Host {
	case "host1":
		targetURL, _ = url.Parse("http://localhost:8080")
	case "host2":
		targetURL, _ = url.Parse("http://localhost:8081")
	default:
		c.String(http.StatusBadGateway, "Invalid host")
		return
	}

	// Create a new reverse proxy for the target URL
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Edit request header
	c.Request.Header.Add("X-Custom-Header", "Custom value")

	// Use reverse proxy
	proxy.ServeHTTP(c.Writer, c.Request)
}
