package consumer

import (
	"context"
	"log"

	"processor/config"
	"processor/worker"

	"github.com/IBM/sarama"
)

func StartConsumer(topic string, workerCount int) {

	cfg := config.NewKafkaConfig()
	cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

	brokers := config.GetBrokers()
	groupID := "log-processor-group"

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, cfg)
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}
	defer consumerGroup.Close()

	jobs := make(chan []byte, 100)

	// start workers
	for i := 0; i < workerCount; i++ {
		go worker.StartWorker(i, jobs)
	}

	handler := &ConsumerHandler{
		jobs: jobs,
	}

	for {
		err := consumerGroup.Consume(context.Background(), []string{topic}, handler)
		if err != nil {
			log.Printf("Error consuming: %v", err)
		}
	}
}