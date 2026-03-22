package pipeline

import (
	"encoding/json"
	"fmt"
	"strings"

	"processor/alert"
	"processor/model"
)

func ProcessLog(data []byte) {

	logData := parseLog(data)

	enriched := enrichLog(logData)

	analyzeLog(enriched)

	checkAlert(enriched)
}

func parseLog(data []byte) model.Log {
	var logData model.Log

	err := json.Unmarshal(data, &logData)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return model.Log{}
	}

	return logData
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

	if log.Level == "ERROR" {
		alert.TriggerAlert(log)
		return
	}

	if strings.Contains(strings.ToLower(log.Message), "timeout") {
		alert.TriggerAlert(log)
		return
	}
}