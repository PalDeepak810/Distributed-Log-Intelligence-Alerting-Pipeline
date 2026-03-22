package consumer

import(
	"log"
	"Processor/config"
	"processor/worker"

	"git.com/IBM/sarama"
)

func StartConsumer(topic string,workerCount int){
	cfg:=config.NewKafkaConfig()
	brokers:=config.GetBrokers()

	consumer, err:=sarama.NewConsumer(brokers,cfg)
	if err!=nil{
		log.Fatalf("Error creating consumer:%v",err)
	}
	defer consumer.Close()

	partitions, err :=consumer.Partitions(topic)
	if err!=nil{
		log.Fatalf("Error fetching partitions: %v",err)
	}

	jobs:=make(chan []byte,100)

	//start workers
	for i:=0,i<workerCount;i++{
		go worker.StartWorker(i,jobs)
	}

	//consume partitions

	for_, partition:=range partitions{
		pc, err:=consumer.ConsumePartition(topic,partition,sarama.OffsetNewest)
		if err!=nil{
			log.Fatalf("Error  consuming partition %d: %v",partition,err)
		}

		go func(pc sarama.PartitionConsumer){
			for:=msg:=range pc.Messages(){
				jobs<-msg.value
			}
		}(pc)
	}
}