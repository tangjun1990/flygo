package kkafka

import (
	"sync"

	"github.com/Shopify/sarama"
	"github.com/tangjun1990/flygo/core/klog"
)

const PackageName = "component.kkafka"

type Component struct {
	config *config
	logger *klog.Component

	consumers     map[string]sarama.Consumer
	producers     map[string]sarama.SyncProducer
	consumerMutex sync.RWMutex
	producerMutex sync.RWMutex
}

func (c *Component) Consumer(name string) sarama.Consumer {
	c.consumerMutex.RLock()
	if consumer, ok := c.consumers[name]; ok {
		c.consumerMutex.RUnlock()
		return consumer
	}
	c.consumerMutex.RUnlock()

	c.consumerMutex.Lock()
	if consumer, ok := c.consumers[name]; ok {
		c.consumerMutex.Unlock()
		return consumer
	}
	_, ok := c.config.Consumers[name]
	if !ok {
		c.logger.Panic("consumer config not exists")
	}

	saramaConfig := sarama.NewConfig()
	saramaConfig.Net.SASL.Enable = true
	saramaConfig.Net.SASL.User = ""
	saramaConfig.Net.SASL.Password = ""
	saramaConfig.Version = sarama.V2_3_1_0
	//saramaConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{}

	consumer, err := sarama.NewConsumer(c.config.Brokers, saramaConfig)
	if err != nil {
		c.logger.Panicf("new consumer err: %v", err.Error())
	}
	c.consumers[name] = consumer
	c.consumerMutex.Unlock()
	return c.consumers[name]
}

func (c *Component) Producer(name string) sarama.SyncProducer {
	c.producerMutex.RLock()
	if producer, ok := c.producers[name]; ok {
		c.producerMutex.RUnlock()
		return producer
	}
	c.producerMutex.RUnlock()

	c.producerMutex.Lock()
	if producer, ok := c.producers[name]; ok {
		c.producerMutex.Unlock()
		return producer
	}
	_, ok := c.config.Producers[name]
	if !ok {
		c.logger.Panic("consumer config not exists")
	}

	saramaConfig := sarama.NewConfig()
	saramaConfig.Net.SASL.Enable = true
	saramaConfig.Net.SASL.User = ""
	saramaConfig.Net.SASL.Password = ""
	saramaConfig.Version = sarama.V2_3_1_0

	producer, err := sarama.NewSyncProducer(c.config.Brokers, saramaConfig)
	if err != nil {
		c.logger.Panicf("new syncProducer err: %v", err.Error())
	}
	c.producers[name] = producer
	c.producerMutex.Unlock()
	return c.producers[name]
}
