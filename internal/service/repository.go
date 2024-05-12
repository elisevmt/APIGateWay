package service

type Repository interface {
	GetServiceProxyInstanceById(serviceId *int64) (*ServiceProxyInstance, error)
	GetServiceProxyInstance() ([]*ServiceProxyInstance, error)
}
