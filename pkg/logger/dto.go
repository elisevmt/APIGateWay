package logger

type LogParams struct {
	ServiceName string `json:"service_name"`
	LogLevel    int64  `json:"log_level"`
	Message     string `json:"message"`
}
