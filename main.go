package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var version = "1.9"

func main() {
	// The "HandleFunc" method accepts a path and a function as arguments
	// (Yes, we can pass functions as arguments, and even trat them like variables in Go)
	// However, the handler function has to have the appropriate signature (as described by the "handler" function below)
	http.HandleFunc("/", handler)
	http.HandleFunc("/env", envHandler)

	// After defining our server, we finally "listen and serve" on port 8080
	// The second argument is the handler, which we will come to later on, but for now it is left as nil,
	// and the handler defined above (in "HandleFunc") is used
	http.ListenAndServe(":8080", nil)
}

// "handler" is our handler function. It has to follow the function signature of a ResponseWriter and Request type
// as the arguments.
func handler(w http.ResponseWriter, r *http.Request) {
	// For this case, we will always pipe "Hello World" into the response writer
	json := "{\"version\":\"" + version + "\"}"
	errorCodes, errorExists := r.URL.Query()["errorcode"]

	if errorExists {
		if errorCode, err := strconv.Atoi(errorCodes[0]); err == nil {
			w.WriteHeader(errorCode)
		}
	}
	fmt.Fprintf(w, json)

}

func envHandler(w http.ResponseWriter, r *http.Request) {
	// For this case, we will always pipe "Hello World" into the response writer
	json := "{\"version\":\"" + version + "\","
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		json = json + "\"" + pair[0] + "\":" + "\"" + pair[1] + "\","
	}
	json = strings.TrimRight(json, ",") + "}"
	fmt.Fprintf(w, json)
	if os.Getenv("DOWNSTREAM") != "" {
		fmt.Fprintf(w, getDownstream(os.Getenv("DOWNSTREAM")))
	}
}
func getDownstream(url string) string {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	response, err := netClient.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	return string(body)
}
