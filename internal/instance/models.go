package instance

type Instance struct {
	Id        *int64  `json:"id" db:"id"`
	ServiceId *int64  `json:"service_id" db:"service_id"`
	Url       *string `json:"url" db:"url"`
	Active    bool    `json:"active" db:"active"`
}
