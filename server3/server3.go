package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var address = ":9002"
var baseURL = "http://localhost"

type FlightData struct {
	DistributionType string `json:"distribution_type"`
	SupplierId       string `json:"supplier_id"`
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "apa kabar!")
}

func handle_hit_get(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get hit from GET")
	fmt.Println(r.Header)
	fmt.Println(r.URL)
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		var result, err = json.Marshal(FlightData{
			DistributionType: "dummy_distribution_type",
			SupplierId:       "dummy_supplier_id",
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(result)
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func handle_hit_post(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get hit from POST")
	fmt.Println(r.Header)
	fmt.Println(r.URL)
	fmt.Println(r.Method)
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("error read response body", err)
		}
		fmt.Println("Body Bytes =\n", string(bodyBytes))
		result := &FlightData{}

		err = json.Unmarshal(bodyBytes, result)

		if err != nil {
			fmt.Println("Error unmarshalling", err)
		}

		fmt.Println(result)

		w.Write(bodyBytes)
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Get hit from server proxy")
		fmt.Println(r.Header)
		fmt.Println(r.URL)
		fmt.Fprintln(w, "port "+address+"testing estu")
	})

	http.HandleFunc("/index", index)

	http.HandleFunc("/hit-server-get", handle_hit_get)
	http.HandleFunc("/hit-server-post", handle_hit_post)

	fmt.Printf("server started at %s%s\n", baseURL, address)

	server := new(http.Server)
	server.Addr = address
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err.Error())
	}
}
