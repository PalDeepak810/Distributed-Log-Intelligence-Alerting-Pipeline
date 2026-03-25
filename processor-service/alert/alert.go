package alert

import (
	"fmt"
	"strings"

	"processor/model"
	"processor/monitor"
)

func TriggerAlert(log model.Log) {
	fmt.Printf(" ALERT: [%s] %s - %s\n",
		log.Service,
		log.Level,
		log.Message,
	)

	alertType := "ALERT"
	if strings.ToUpper(strings.TrimSpace(log.Level)) == "ERROR" {
		alertType = "ERROR_LOG"
	}

	monitor.AddAlert(monitor.AlertEvent{
		Type:    alertType,
		Service: log.Service,
		Level:   log.Level,
		Message: log.Message,
	})
}
