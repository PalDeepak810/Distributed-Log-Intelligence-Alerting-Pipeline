package pipeline

import (
	"encoding/json"
	"fmt"
	"strings"

	"processor/alert"
	"processor/model"
)

func ProcessLog(data []byte) {

	var raw map[string]interface{}

	err := json.Unmarshal(data, &raw)
	if err != nil {
		fmt.Println("Invalid JSON:", err)
		return
	}

	//  Normalize here
	logData := Normalize(raw)

	enriched := enrichLog(logData)

	analyzeLog(enriched)

	checkAlert(enriched)
}

func enrichLog(log model.Log) model.Log {
	if log.Level == "" {
		log.Level = "INFO"
	}
	return log
}

func analyzeLog(log model.Log) {
	fmt.Println("Processed Log:", log.Service, log.Level, log.Message)
}

func checkAlert(log model.Log) {
	level := strings.ToUpper(strings.TrimSpace(log.Level))
	timeoutDetected := strings.Contains(strings.ToLower(log.Message), "timeout")
	errorDetected := level == "ERROR"

	if errorDetected || timeoutDetected {
		alert.TriggerAlert(log)
		alert.CheckErrorSpike(log.Service)
	}
}
