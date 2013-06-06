package main

import (
	// "encoding/json"
	// "flag"
	// "fmt"
	"github.com/alexgoodell/ghdmodel/models"
	// "io/ioutil"
	// "net/http"
	// "strconv"
)

func main() {

	costanalysis.Predict()

	// port := flag.Int("port", 3000, "server port")

	// fmt.Println("Starting webserver. Listenning on port", *port)
	// http.HandleFunc("/cost_analysis", costAnalysisHandler)
	// http.ListenAndServe(":"+strconv.Itoa(*port), nil)
}

// func costAnalysisHandler(respWriter http.ResponseWriter, req *http.Request) {
// 	fmt.Println("GET /cost_analysis")
// 	body, err := ioutil.ReadAll(req.Body)
// 	if err != nil {
// 		panic("Cannot read body?")
// 	}
// 	inputs := new(costanalysis.Inputs)
// 	json.Unmarshal(body, inputs)
// 	if err != nil {
// 		panic("Json error")
// 	}
// 	results := costanalysis.Predict(inputs)
// 	responseBody, err := json.Marshal(results)
// 	respWriter.Header().Set("Content-Type", "application/json")
// 	if err == nil {
// 		fmt.Fprintf(respWriter, string(responseBody))
// 	} else {
// 		http.Error(respWriter, err.Error(), http.StatusInternalServerError)
// 	}
// }
