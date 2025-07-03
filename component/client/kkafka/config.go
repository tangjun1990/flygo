package kkafka

type config struct {
	Brokers []string
	Debug   bool

	Producers map[string]producerConfig // 目前仅支持同步生产者，异步生产者暂不支持
	Consumers map[string]consumerConfig
}

const (
	balancerHash = "hash"
	balancerRR   = "roundRobin"
)

func DefaultConfig() *config {
	return &config{}
}

type producerConfig struct {
	Topic    string
	Balancer string
	Async    bool
}

type consumerConfig struct {
	Topic     string
	Partition int
	GroupID   string
}
