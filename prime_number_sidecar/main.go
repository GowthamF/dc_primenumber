package main

import (
	"bufio"
	"flag"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"dc_assignment.com/sidecar/v2/sidecarroutes"
	"github.com/gin-gonic/gin"
)

var (
	appPortNumber     = flag.Int("appPortNumber", 8080, "App Port Number")
	sideCarPortNumber = flag.Int("sideCarPortNumber", 8081, "Sidecar Port Number")
)

func main() {
	r := sidecarroutes.SetupRouter(strconv.Itoa(*appPortNumber))
	r.GET("/:path", func(c *gin.Context) {
		// step 1: resolve proxy address, change scheme and host in requets
		req := c.Request
		proxy, err := url.Parse("http://localhost:" + strconv.Itoa(*appPortNumber))
		if err != nil {
			log.Printf("error in parse addr: %v", err)
			c.String(500, "error")
			return
		}
		req.URL.Scheme = proxy.Scheme
		req.URL.Host = proxy.Host

		// step 2: use http.Transport to do request to real server.
		transport := http.DefaultTransport
		resp, err := transport.RoundTrip(req)
		if err != nil {
			log.Printf("error in roundtrip: %v", err)
			c.String(500, "error")
			return
		}

		// step 3: return real server response to upstream.
		for k, vv := range resp.Header {
			for _, v := range vv {
				c.Header(k, v)
			}
		}
		defer resp.Body.Close()
		bufio.NewReader(resp.Body).WriteTo(c.Writer)
		return
	})
	r.Run(":" + strconv.Itoa(*sideCarPortNumber))
}
