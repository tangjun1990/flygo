package main

import (
	"fmt"
	"github.com/Shopify/sarama"
)

func main() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	p, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		fmt.Println("NewSyncProducer err", err.Error())
	}
	defer p.Close()

	msgs := make([]*sarama.ProducerMessage, 0)
	for i := 0; i < 100; i++ {
		msg := &sarama.ProducerMessage{}
		msg.Topic = "demosms"
		msg.Value = sarama.StringEncoder(fmt.Sprintf("sajdaskjdaskj %d", i))

		msgs = append(msgs, msg)
	}
	err = p.SendMessages(msgs)
	fmt.Println("sendMessage", err)
}
