package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ximus/ghimodel/models"
	"io/ioutil"
	"net/http"
	"strconv"
)

func main() {
	port := flag.Int("port", 3000, "server port")

	fmt.Println("Starting webserver. Listenning on port", *port)
	http.HandleFunc("/cost_analysis", costAnalysisHandler)
	http.ListenAndServe(":"+strconv.Itoa(*port), nil)
}

func costAnalysisHandler(respWriter http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic("Cannot read body?")
	}
	inputs := &costanalysis.Inputs{}
	json.Unmarshal(body, &inputs)
	if err != nil {
		panic("Json error")
	}
	results := costanalysis.Predict(inputs)
	responseBody, err := json.Marshal(results)
	respWriter.Header().Set("Content-Type", "application/json")
	if err == nil {
		fmt.Fprintf(respWriter, string(responseBody))
	} else {
		http.Error(respWriter, err.Error(), http.StatusInternalServerError)
	}
}
