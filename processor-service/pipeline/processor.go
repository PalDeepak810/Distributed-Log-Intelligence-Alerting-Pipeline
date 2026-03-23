package pipeline

import(
	"encoding/json"
	"fmt"
	"processor/model"
)

func ProcessLog(data []byte){
	logData:=parseLog(data)

	enriched:=enrichLog(logData)

	analyze

}