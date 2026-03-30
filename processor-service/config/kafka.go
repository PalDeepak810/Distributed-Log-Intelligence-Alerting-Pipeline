package config

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/IBM/sarama"
)

func NewKafkaConfig() *sarama.Config {
	config := sarama.NewConfig()

	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	//Kafka version
	config.Version = sarama.V2_1_0_0

	//enable TLS (AIVEN requirement)
	config.Net.TLS.Enable = true

	//Load CA cert
	caCert, err := ioutil.ReadFile("ca.pem")
	if err != nil {
		panic(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	//Load client cert+key
	cert, err := tls.LoadX509KeyPair("service.cert", "service.key")
	if err != nil {
		panic(err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	config.Net.TLS.Config = tlsConfig

	return config
}

func GetBrokers() []string {
	return []string{
		"kafkalogservice-dp220001-beec.g.aivencloud.com:26932",
	}
}
