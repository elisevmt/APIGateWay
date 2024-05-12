package proxy

type Proxy struct {
	Id       *int64  `json:"id" db:"id"`
	Strategy *string `json:"strategy" db:"strategy"`
}
