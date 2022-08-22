package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type FlightData struct {
	DistributionType string `json:"distribution_type"`
	SupplierId       string `json:"supplier_id"`
}

func hitServer2(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hitting Server... = ", r.URL)

	var request *http.Request
	var err error
	if r.Method == "GET" {
		fmt.Println("hit /hit-server-get")
		request, err = http.NewRequest("GET", "http://localhost:9001/hit-server-get", nil)
	} else if r.Method == "POST" {
		fmt.Println("hit /hit-server-post")
		request, err = http.NewRequest("POST", "http://localhost:9001/hit-server-post", r.Body)
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	request.Header.Set("come-from", "server-1")

	var client = &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()
	fmt.Println("Response Header = ", response.Header)

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("error read response body", err)
	}
	fmt.Println("Incoming Body Bytes = ", string(bodyBytes))
	result := &FlightData{}

	err = json.Unmarshal(bodyBytes, result)

	if err != nil {
		fmt.Println("Error unmarshalling", err)
	}

	fmt.Println(result)
}

func hitServer3(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hitting Server... = ", r.URL)

	var request *http.Request
	var err error
	if r.Method == "GET" {
		fmt.Println("hit /hit-server-get")
		request, err = http.NewRequest("GET", "http://localhost:9002/hit-server-get", nil)
	} else if r.Method == "POST" {
		fmt.Println("hit /hit-server-post")
		request, err = http.NewRequest("POST", "http://localhost:9002/hit-server-post", r.Body)
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	var client = &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()
	fmt.Println("Response Header = ", response.Header)

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("error read response body", err)
	}
	fmt.Println("Body Bytes = ", string(bodyBytes))
	result := &FlightData{}

	err = json.Unmarshal(bodyBytes, result)

	if err != nil {
		fmt.Println("Error unmarshalling", err)
	}

	fmt.Println(result)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "This is Server 1 on port 9000")
	})

	http.HandleFunc("/hit-server-2", hitServer2)
	http.HandleFunc("/hit-server-3", hitServer3)

	fmt.Println("server started at http://localhost:9000")
	http.ListenAndServe("localhost:9000", nil)
}
