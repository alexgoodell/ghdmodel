
package main

import (

	"fmt"
	"costanalysis"

)



// type App struct {
// 	env string
// 	logger log.Logger
// }

// func NewApp() {
// 	env    := flag.String("env", "development", "App environment")
// 	logger := newLogger()

// 	app := App{env: env, logger: logger}
// }

// var app = NewApp()

// func main() {
	// port := flag.Int("port", 3000, "server port")
	// http.HandleFunc("/cost_analysis", costAnalysisHandler)
	// http.ListenAndServe(":"+port, nil)
// }

// func costAnalysisHandler(w http.ResponseWriter, r *http.Request) {
// 	var inputs = costanalysis.Inputs
// 	var results = costanalysis.Predict(&inputs)
// 	var responseBody, err = json.Marshal(results)
// 	if err == nil {
// 		fmt.Fprintf(w, responseBody, r.URL.Path[1:])
// 	} else {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

// func newLogger(path string, level int) fmt.Logger {
// 	w := bufio.NewWriter(os.Stdout)
// 	logger := log.New(w)
// 	return &logger;
// }


func main() {


	Predict()
	fmt.Println("Done")


}

