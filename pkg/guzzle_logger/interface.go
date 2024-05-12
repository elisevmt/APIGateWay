package guzzle_logger

type API interface {
	SendLog(level, msg, desciption string) error
}
