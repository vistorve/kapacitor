package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/influxdata/kapacitor"
	"log"
)

type Service struct {
	topics           []string
	global           bool
	stateChangesOnly bool
	producer         sarama.SyncProducer
	logger           *log.Logger
}

func NewService(c Config, l *log.Logger) *Service {
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Retry.Max = 10
	// SASL
	conf.Net.SASL.Enable = c.SaslEnabled
	conf.Net.SASL.User = c.SaslUser
	conf.Net.SASL.Password = c.SaslPassword
	producer, err := sarama.NewSyncProducer(c.Brokers, conf)

	if err != nil {
		log.Fatalf("Cannot create kafka producer: %v", err)
	}

	return &Service{
		topics:           c.Topics,
		global:           c.Global,
		stateChangesOnly: c.StateChangesOnly,
		producer:         producer,
		logger:           l,
	}
}

func NewMockService(c Config, l *log.Logger, p sarama.SyncProducer) *Service {
	return &Service{
		topics:           c.Topics,
		global:           c.Global,
		stateChangesOnly: c.StateChangesOnly,
		producer:         p,
		logger:           l,
	}
}

func (s *Service) Open() error {
	return nil
}

func (s *Service) Close() error {
	s.producer.Close()
	return nil
}

func (s *Service) Global() bool {
	return s.global
}
func (s *Service) StateChangesOnly() bool {
	return s.stateChangesOnly
}

func (s *Service) Alert(topics []string, message string, level kapacitor.AlertLevel) error {
	if len(topics) == 0 {
		topics = s.topics
	}
	for _, topic := range topics {
		_, _, err := s.producer.SendMessage(&sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(message)})

		if err != nil {
			return err
		}
	}

	return nil
}