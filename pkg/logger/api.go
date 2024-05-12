package logger

import (
	"APIGateWay/pkg/rabbit/publisher"
	"encoding/json"
	"fmt"
	"runtime"
)

type loggerAPI struct {
	appName         string
	rabbitPublisher publisher.Publisher
}

func NewLoggerAPI(rabbitUrl, logQueue, appName string) API {
	rabbitPublisher, err := publisher.NewPublisher(rabbitUrl, logQueue)
	if err != nil {
		panic(err)
	}
	return &loggerAPI{
		rabbitPublisher: *rabbitPublisher,
		appName:         appName,
	}
}

func (l *loggerAPI) ErrorLog(error error) {
	_, fn, line, _ := runtime.Caller(1)
	msg := fmt.Sprintf("ERROR:\n%s :: %d :: %s", fn, line, error.Error())
	logParams := LogParams{
		ServiceName: l.appName,
		LogLevel:    ERROR_LVL,
		Message:     msg,
	}
	buffer, _ := json.Marshal(logParams)
	l.rabbitPublisher.Publish(buffer)
}

func (l *loggerAPI) MsgLog(msg string) {
	logParams := LogParams{
		ServiceName: l.appName,
		LogLevel:    INFO_LVL,
		Message:     msg,
	}
	buffer, _ := json.Marshal(logParams)
	l.rabbitPublisher.Publish(buffer)
}
