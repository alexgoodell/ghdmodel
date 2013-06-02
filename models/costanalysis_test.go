package costanalysis

import (
	"encoding/json"
	"testing"
)

const (
	inputsPath = "test/sample_inputs.json"
	resultPath = "test/sample_results.json"
)

func TestPredict(t *testing.T) {
	queryJson, _ := ioutil.ReadFile(queryPath)
	inputs := new(Inputs)
	json.Unmarshal(queryJson, inputs)

	resultsJson, _ := ioutil.ReadFile(resultPath)
	results := new(Results)
	json.Unmarshal(resultsJson, results)

	if (*results == *Predict(inputs)) {
		t.Errorf("Sqrt(%v) = %v, want %v", in, x, out)
	}
}