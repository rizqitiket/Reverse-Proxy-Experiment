package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var address = ":9001"
var baseURL = "http://localhost"
var limiter = 5
var key_limiter = map[string]int{}

type FlightData struct {
	DistributionType string `json:"distribution_type"`
	SupplierId       string `json:"supplier_id"`
}

// NewProxy takes target host and creates a reverse proxy
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	fmt.Println("URL = ", url)
	proxy.ModifyResponse = modifyResponse()
	return proxy, nil
}

func modifyResponse() func(*http.Response) error {
	return func(resp *http.Response) error {
		resp.Header.Set("X-Proxy", "B to A")
		return nil
	}
}

// ProxyRequestHandler handles the http request using proxy
func ProxyRequestHandler(proxy1 *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Get hit from", r.Host, r.URL)
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("error read response body", err)
		}
		r.Body.Close()
		//to reproduce request body after close it
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		fmt.Println("Incoming Header from server 1 =\n", r.Header)
		fmt.Println("Incoming Body Bytes from server 1 =\n", string(bodyBytes))

		result := &FlightData{}

		err = json.Unmarshal(bodyBytes, result)

		if err != nil {
			fmt.Println("Error unmarshalling", err)
		}

		key_limiter[result.DistributionType] += 1
		fmt.Println("Limiter = ", key_limiter)

		if key_limiter[result.DistributionType] >= 7 {
			key_limiter[result.DistributionType] = 0
		}

		if key_limiter[result.DistributionType] < limiter {
			proxy1.ServeHTTP(w, r)
		} else {
			fmt.Fprintln(w, "booking limit has been exceeded")
		}
	}
}

func main() {
	// initialize a reverse proxy and pass the actual backend server url here
	proxy1, err := NewProxy(baseURL + ":9002")
	if err != nil {
		panic(err)
	}

	// handle all requests to your server using the proxy
	http.HandleFunc("/", ProxyRequestHandler(proxy1))
	fmt.Printf("server started at %s%s\n", baseURL, address)
	http.ListenAndServe(address, nil)
}
