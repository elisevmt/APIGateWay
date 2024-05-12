package mutex_tool

import (
	"sync"
)

type Tool struct {
	createMute     *sync.Mutex
	clientMutexMap map[string]map[string]*sync.Mutex
}

func NewTool() *Tool {
	return &Tool{createMute: new(sync.Mutex), clientMutexMap: make(map[string]map[string]*sync.Mutex)}
}

func (t *Tool) RegisterClientMutex(clientId string, targetKey string) {
	t.createMute.Lock()
	if _, ok := t.clientMutexMap[clientId]; ok {
		t.clientMutexMap[clientId][targetKey] = new(sync.Mutex)
	} else {
		t.clientMutexMap[clientId] = make(map[string]*sync.Mutex)
		t.clientMutexMap[clientId][targetKey] = new(sync.Mutex)
	}
	t.createMute.Unlock()
}

func (t *Tool) Fire(clientId string, targetKey string) {
	if t.clientMutexMap[clientId][targetKey] == nil {
		t.RegisterClientMutex(clientId, targetKey)
	}
}

func (t *Tool) LockClient(clientId string, targetKey string) {
	t.clientMutexMap[clientId][targetKey].Lock()
}

func (t *Tool) UnlockClient(clientId string, targetKey string) {
	t.clientMutexMap[clientId][targetKey].Unlock()
}
