package consumer

import (
	"github.com/wagslane/go-rabbitmq"
)

func StartConsuming(
	url, queueName string,
	consumeFunc func(d []byte) error) error {
	conn, err := rabbitmq.NewConn(url)
	if err != nil {
		return err
	}
	_, err = rabbitmq.NewConsumer(
		conn,
		func(d rabbitmq.Delivery) rabbitmq.Action {
			consumeFunc(d.Body)
			return rabbitmq.Ack
		},
		queueName,
		rabbitmq.WithConsumerOptionsLogging,
	)
	if err != nil {
		return err
	}
	return nil
}
