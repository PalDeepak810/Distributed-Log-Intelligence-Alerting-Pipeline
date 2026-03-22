package pipeline

import(
	"encoding/json"
	"fmt"
	"processor/model"
)

func ProcessLog(data []byte){
	logData:=parseLog(data)
	enriched:=enrichLog(logData)
	analyzeLog(enriched)

}

//-Internal steps--

func parseLog(data []byte) model.Log{
	var logData model.Log

	err:=json.Unmarshal(data,&logData)
	if err!=nil{
		fmt.Println("error parsing JSOn:",err)
		return model.Log{}
	}

	return logData
}

func enrichLog(log model.Log){
	fmt.println("Processed Log:",log.Service,log.Level,log.Message)
}