package alert

import (
	"fmt"
	"processor/model"
)

func TriggerAlert(log model.Log) {
	fmt.Printf(" ALERT: [%s] %s - %s\n",
		log.Service,
		log.Level,
		log.Message,
	)
}