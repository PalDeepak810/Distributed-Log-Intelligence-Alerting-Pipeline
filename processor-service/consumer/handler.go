package consumer

import (
	"github.com/IBM/sarama"
)

type ConsumerHandler struct {
	jobs chan []byte
}

// called at start
func (h *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// called at end
func (h *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// main logic
func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for msg := range claim.Messages() {
		h.jobs <- msg.Value

		// mark message as processed
		session.MarkMessage(msg, "")
	}

	return nil
}