package proxy_uc

import (
	"APIGateWay/internal/instance"
	"APIGateWay/internal/proxy"
	"APIGateWay/internal/service"
	GL "APIGateWay/pkg/guzzle_logger"
	"APIGateWay/pkg/mutex_tool"
)

type proxyUC struct {
	mutextTool                   *mutex_tool.Tool
	serviceInstanceMap           map[int64][]*instance.Instance
	serviceStategyMap            map[int64]string
	instanceLoadMap              map[string]int64
	instanceTotalRequestCountMap map[string]int64
	serviceTotalRequestCountMap  map[int64]int64
	serviceRepo                  service.Repository
	gl                           GL.API
}

func NewProxyUC(serviceRepo service.Repository, gl GL.API) proxy.UC {
	serviceInstanceMap := make(map[int64][]*instance.Instance)
	serviceStategyMap := make(map[int64]string)
	instanceLoadMap := make(map[string]int64)
	serviceTotalRequestCountMap := make(map[int64]int64)
	instanceTotalRequestCountMap := make(map[string]int64)
	serviceProxyInstances, err := serviceRepo.GetServiceProxyInstance()
	if err != nil {
		panic(err)
	}
	for _, spi := range serviceProxyInstances {
		if !spi.Active {
			continue
		}
		serviceInstanceMap[*spi.ServiceId] = append(serviceInstanceMap[*spi.ServiceId], &spi.Instance)
		serviceStategyMap[*spi.ServiceId] = *spi.Strategy
		instanceLoadMap[*spi.Url] = 0
		instanceTotalRequestCountMap[*spi.Url] = 0
		serviceTotalRequestCountMap[*spi.ServiceId] = 0
	}
	tool := mutex_tool.NewTool()
	tool.Fire("serviceInstanceMap", "0")
	tool.Fire("serviceTotalRequestCountMap", "0")
	tool.Fire("serviceStategyMap", "0")
	tool.Fire("instanceLoadMap", "0")
	tool.Fire("instanceTotalRequestCountMap", "0")
	gl.SendLog(GL.LEVEL_INFO, "-", "System startup")
	return &proxyUC{
		serviceInstanceMap:           serviceInstanceMap,
		serviceTotalRequestCountMap:  serviceTotalRequestCountMap,
		serviceStategyMap:            serviceStategyMap,
		instanceLoadMap:              instanceLoadMap,
		instanceTotalRequestCountMap: instanceTotalRequestCountMap,
		serviceRepo:                  serviceRepo,
		mutextTool:                   tool,
		gl:                           gl,
	}
}

func (p *proxyUC) revolver(serviceId *int64) (*instance.Instance, error) {

	p.mutextTool.LockClient("serviceInstanceMap", "0")
	instanceListLength := len(p.serviceInstanceMap[*serviceId])
	p.mutextTool.UnlockClient("serviceInstanceMap", "0")

	p.mutextTool.LockClient("serviceTotalRequestCountMap", "0")
	index := p.serviceTotalRequestCountMap[*serviceId] % int64(instanceListLength)
	p.mutextTool.UnlockClient("serviceTotalRequestCountMap", "0")

	p.mutextTool.LockClient("serviceTotalRequestCountMap", "0")
	p.serviceTotalRequestCountMap[*serviceId] += 1
	p.mutextTool.UnlockClient("serviceTotalRequestCountMap", "0")

	p.mutextTool.LockClient("instanceLoadMap", "0")
	p.instanceLoadMap[*p.serviceInstanceMap[*serviceId][index].Url] += 1
	p.mutextTool.UnlockClient("instanceLoadMap", "0")

	p.mutextTool.LockClient("instanceTotalRequestCountMap", "0")
	p.instanceTotalRequestCountMap[*p.serviceInstanceMap[*serviceId][index].Url] += 1
	p.mutextTool.UnlockClient("instanceTotalRequestCountMap", "0")
	return p.serviceInstanceMap[*serviceId][index], nil
}

func (p *proxyUC) activeLoad(serviceId *int64) (*instance.Instance, error) {
	var targetInstance *instance.Instance
	minimalLoad := 10000000000

	p.mutextTool.LockClient("serviceInstanceMap", "0")
	instances := p.serviceInstanceMap[*serviceId]
	p.mutextTool.UnlockClient("serviceInstanceMap", "0")
	for _, instance := range instances {
		p.mutextTool.LockClient("instanceLoadMap", "0")
		if p.instanceLoadMap[*instance.Url] < int64(minimalLoad) {
			minimalLoad = int(p.instanceLoadMap[*instance.Url])
			targetInstance = instance
		}
		p.mutextTool.UnlockClient("instanceLoadMap", "0")
	}
	p.mutextTool.LockClient("instanceLoadMap", "0")
	p.instanceLoadMap[*targetInstance.Url] += 1
	p.mutextTool.UnlockClient("instanceLoadMap", "0")
	return targetInstance, nil
}

func (p *proxyUC) DecreaseLoad(instanceUrl *string) {
	p.mutextTool.LockClient("instanceLoadMap", "0")
	p.instanceLoadMap[*instanceUrl] -= 1
	p.mutextTool.UnlockClient("instanceLoadMap", "0")
}

func (p *proxyUC) GetProxyInstance(serviceId *int64) (*instance.Instance, error) {
	switch p.serviceStategyMap[*serviceId] {
	case "revolver":
		return p.revolver(serviceId)
	case "active_load":
		return p.activeLoad(serviceId)
	default:
		return p.revolver(serviceId)
	}
}
