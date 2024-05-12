package crud

type CRUD interface {
	Create(data interface{}) error
	Get(filters map[string]interface{}) (interface{}, error)
	Fetch(filters map[string]interface{}, order *string, orderType *string, limit *int64) ([]interface{}, error)
	Update(filters map[string]interface{}, units map[string]interface{}) error
	Delete(filters map[string]interface{}) error
}
