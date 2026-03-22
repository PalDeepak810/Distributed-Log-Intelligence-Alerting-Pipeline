package config
import "github.com/IBM/sarama"

func NewKafkaConfig() *sarama.Config{
	config:=sarama.NewConfig()

	config.Consume.return.Errors=true
	config.Consumer.Offsets.Initial=sarama.OffsetNewest
	config.version=sarama.V2_1_0_0

	return config
}

func GetBrokers() []string{
	return []string{"127.0.0.1:9092"}
}