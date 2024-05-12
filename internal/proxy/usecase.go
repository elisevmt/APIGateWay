package proxy

import "APIGateWay/internal/instance"

type UC interface {
	GetProxyInstance(serviceId *int64) (*instance.Instance, error)
	DecreaseLoad(instanceUrl *string)
}
