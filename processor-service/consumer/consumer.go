package consumer

import (
	"log"

	"processor/config"
	"processor/worker"

	"github.com/IBM/sarama"
)

func StartConsumer(topic string, workerCount int) {

	cfg := config.NewKafkaConfig()
	brokers := config.GetBrokers()

	consumer, err := sarama.NewConsumer(brokers, cfg)
	if err != nil {
		log.Fatalf("Error creating consumer: %v", err)
	}

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Fatalf("Error fetching partitions: %v", err)
	}
	log.Printf("Kafka consumer connected. topic=%s partitions=%v", topic, partitions)

	jobs := make(chan []byte, 100)

	// start workers
	for i := 0; i < workerCount; i++ {
		go worker.StartWorker(i, jobs)
	}

	// consume partitions
	for _, partition := range partitions {

		pc, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Fatalf("Error consuming partition %d: %v", partition, err)
		}

		go func(partition int32, pc sarama.PartitionConsumer) {
			for {
				select {
				case msg, ok := <-pc.Messages():
					if !ok {
						log.Printf("Partition %d message channel closed", partition)
						return
					}
					jobs <- msg.Value
				case err, ok := <-pc.Errors():
					if ok && err != nil {
						log.Printf("Partition %d consumer error: %v", partition, err)
					}
				}
			}
		}(partition, pc)
	}
}
