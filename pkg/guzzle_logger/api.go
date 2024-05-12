package guzzle_logger

import (
	"APIGateWay/pkg/rabbit/publisher"
	"encoding/json"
	"fmt"
)

type GuzzleAPI struct {
	remoteService    *string
	remoteSubService *string
	logPublisher     *publisher.Publisher
}

func New(service, subService string, logPublisher *publisher.Publisher) API {
	return &GuzzleAPI{
		remoteService:    &service,
		remoteSubService: &subService,
		logPublisher:     logPublisher,
	}
}

func (g *GuzzleAPI) SendLog(level, msg, description string) error {
	log := Log{
		LogMessage: &LogMessage{
			&msg,
		},
		LogLevel:         &level,
		RemoteService:    g.remoteService,
		RemoteSubService: g.remoteSubService,
		LogDescription:   &description,
	}
	data, err := json.Marshal(log)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	err = g.logPublisher.Publish(data)
	if err != nil {
		return err
	}
	return nil
}
