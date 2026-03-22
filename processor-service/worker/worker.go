package worker

import "processor/pipeline"

func StartWorker(id ,int,jobs<-chan []byte){
	for job:=range job{
		pipeline.ProcessLog(job)
	}
}