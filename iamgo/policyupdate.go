package iamgo

import (
	"sync"

	"cloud.google.com/go/iam/apiv1/iampb"
)

type resourceName string

type policyUpdate struct {
	mutex     *sync.Mutex
	resources map[resourceName]*sync.Mutex
	client    iampb.IAMPolicyClient
}

func newPolicyUpdate(client iampb.IAMPolicyClient) *policyUpdate {
	return &policyUpdate{
		mutex:     &sync.Mutex{},
		resources: make(map[resourceName]*sync.Mutex),
		client:    client,
	}
}

func (s *policyUpdate) lock(resource resourceName) func() {
	s.mutex.Lock()
	mutex, ok := s.resources[resource]
	if !ok {
		mutex = &sync.Mutex{}
		s.resources[resource] = mutex
	}
	s.mutex.Unlock()
	mutex.Lock()
	return mutex.Unlock
}
