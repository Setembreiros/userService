package kafka

import (
	"log"
	"userservice/internal/bus"

	"github.com/IBM/sarama"
)

type Consumer struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	ready    chan bool
	eventBus *bus.EventBus
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				consumer.errorLog.Printf("message channel was closed")
				return nil
			}
			consumer.infoLog.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			session.MarkMessage(message, "")

			event := bus.Event{
				Type: message.Topic,
				Data: message.Value,
			}
			consumer.eventBus.Publish(event)
		case <-session.Context().Done():
			return nil
		}
	}
}
