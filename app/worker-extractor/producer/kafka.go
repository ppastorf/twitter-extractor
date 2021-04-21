package producer

import (
	"fmt"
	"log"
	"os"

	"github.com/Shopify/sarama"
)

type KafkaConfig struct {
	Address    string `yaml:"endpoint"`
	EnableTls  bool   `yaml:"enable_tls"`
	MaxRetries int    `yaml:"max_retries"`
	EnableLogs bool   `yaml:"sarama_logs"`

	Topics struct {
		Workload struct {
			Name       string `yaml:"name"`
			SchemaFile string `yaml:"schema"`
		} `yaml:"workload"`
	} `yaml:"topics"`
}

type Message interface{}

type (
	Producer interface {
		PublishSync(topic string, message []byte) error
		PublishAsync(topic string, message []byte)
		LogEvents()
		Close() error
	}
	producer struct {
		sync  sarama.SyncProducer
		async sarama.AsyncProducer
	}
)

func NewProducer(config *KafkaConfig) (Producer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Net.TLS.Enable = config.EnableTls
	saramaConfig.Producer.Retry.Max = config.MaxRetries
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true
	saramaConfig.Version = sarama.V2_3_0_0
	saramaConfig.Producer.Compression = sarama.CompressionSnappy

	if config.EnableLogs {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	syncHandler, err := sarama.NewSyncProducer([]string{config.Address}, saramaConfig)
	if err != nil {
		return nil, err
	}

	asyncHandler, err := sarama.NewAsyncProducer([]string{config.Address}, saramaConfig)
	if err != nil {
		return nil, err
	}

	prd := &producer{sync: syncHandler, async: asyncHandler}
	return prd, err
}

func (p *producer) PublishSync(topic string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}

	_, _, err := p.sync.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}

func (p *producer) PublishAsync(topic string, message []byte) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}

	p.async.Input() <- msg
}

func (p *producer) LogEvents() {
	for {
		select {
		case msg, ok := <-p.async.Successes():
			if !ok {
				return
			}
			fmt.Printf("Sent to topic %s(p %d) at offset %d\n", msg.Topic, msg.Partition, msg.Offset)
		case msg, ok := <-p.async.Errors():
			if !ok {
				return
			}
			fmt.Printf("Failed to send message: %s\n", msg)
		}
	}
}

func (p *producer) Close() error {
	err := p.sync.Close()
	if err != nil {
		return err
	}

	p.async.Close()
	p.async.AsyncClose()

	return nil
}
