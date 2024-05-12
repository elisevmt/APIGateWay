package logger

type API interface {
	ErrorLog(error error)
	MsgLog(msg string)
}
