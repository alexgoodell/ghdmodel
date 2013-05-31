package main

import (
	"encoding/json"
	"fmt"
	"flag"
	"net/http"
	"strconv"
	"io/ioutil"
	"github.com/ximus/ghimodel/models"
)

func main() {
	port := flag.Int("port", 3000, "server port")

	http.HandleFunc("/cost_analysis", costAnalysisHandler)
	http.ListenAndServe(":"+strconv.Itoa(*port), nil)
}

func costAnalysisHandler(resp http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic("Cannot read body?")
	}
	var inputs costanalysis.Inputs
	json.Unmarshal(body, &inputs)
	if err != nil {
		panic("Json error")
	}
	results := costanalysis.Predict(&inputs)
	responseBody, err := json.Marshal(results)
	resp.Header().Set("Content-Type", "application/json")
	if err == nil {
		fmt.Fprintf(resp, string(responseBody), req.URL.Path[1:])
	} else {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
	}
}
