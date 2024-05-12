package guzzle_logger

const (
	LEVEL_INFO    = "INFO"
	LEVEL_WARNING = "WARNING"
	LEVEL_ERROR   = "ERROR"
)

type Log struct {
	LogMessage       *LogMessage `json:"log_message"`
	LogLevel         *string     `json:"log_level"`
	RemoteService    *string     `json:"remote_service"`
	RemoteSubService *string     `json:"remote_sub_service"`
	LogDescription   *string     `json:"log_description"`
}

type LogMessage struct {
	Body *string `json:"body"`
}
