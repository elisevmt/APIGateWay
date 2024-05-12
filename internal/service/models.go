package service

import (
	"APIGateWay/internal/instance"
	"APIGateWay/internal/proxy"
)

type Service struct {
	Id      *int64  `json:"id" db:"id"`
	Name    *string `json:"name" db:"name"`
	ProxyId *int64  `json:"proxy_id" db:"proxy_id"`
}

type ServiceProxyInstance struct {
	Service
	proxy.Proxy
	instance.Instance
}
