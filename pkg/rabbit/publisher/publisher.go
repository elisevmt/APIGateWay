package publisher

import (
	"github.com/wagslane/go-rabbitmq"
)

type Publisher struct {
	publisher *rabbitmq.Publisher
	queueName string
}

func NewPublisher(url, queueName string) (*Publisher, error) {
	conn, err := rabbitmq.NewConn(url)
	if err != nil {
		return nil, err
	}
	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
	)
	if err != nil {
		return nil, err
	}
	return &Publisher{
		publisher: publisher,
		queueName: queueName,
	}, nil
}

func (p *Publisher) Publish(msg []byte) error {
	err := p.publisher.Publish(
		msg,
		[]string{p.queueName},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsMandatory,
		rabbitmq.WithPublishOptionsPersistentDelivery,
	)
	if err != nil {
		return err
	}
	return nil
}
