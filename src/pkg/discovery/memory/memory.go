package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"moviemicroservice.com/src/pkg/discovery"
)

type ServiceName string

type InstanceID string

type serviceInstance struct {
	hostPort     string
	lastActiveAt time.Time
}

// defines an in-memory registry implementation
type Registry struct {
	*sync.RWMutex
	addresses map[ServiceName]map[InstanceID]*serviceInstance
}

func New() *Registry {
	return &Registry{addresses: map[ServiceName]map[InstanceID]*serviceInstance{}}
}

func (r *Registry) Register(ctx context.Context, instanceID InstanceID, serviceName ServiceName, hostPort string) error {
	r.Lock()

	defer r.Unlock()

	if _, ok := r.addresses[serviceName]; !ok {
		r.addresses[serviceName] = map[InstanceID]*serviceInstance{}
	}

	lastActiveAt := time.Now()
	r.addresses[serviceName][instanceID] = &serviceInstance{hostPort, lastActiveAt}

	return nil
}

func (r *Registry) Deregister(ctx context.Context, instanceID InstanceID, serviceName ServiceName) error {
	r.Lock()

	defer r.Unlock()

	if _, ok := r.addresses[serviceName]; !ok {
		return nil
	}

	delete(r.addresses[serviceName], instanceID)

	return nil
}

func (r *Registry) ReportHealthyState(instanceID InstanceID, serviceName ServiceName) error {
	r.Lock()

	defer r.Unlock()

	if _, ok := r.addresses[serviceName]; !ok {
		return errors.New("service not registered")
	}

	if _, ok := r.addresses[serviceName][instanceID]; !ok {
		return errors.New("service instance not registered")
	}

	r.addresses[serviceName][instanceID].lastActiveAt = time.Now()

	return nil
}

func (r *Registry) ServiceAddresses(ctx context.Context, serviceName ServiceName) ([]string,
	error) {
	r.RLock()

	defer r.RUnlock()

	if len(r.addresses[serviceName]) == 0 {
		return nil, discovery.ErrNotFound
	}

	var res []string

	for _, instance := range r.addresses[serviceName] {
		//skip instances that haven't reported health in the past 5 seconds
		if instance.lastActiveAt.Before(time.Now().Add(-5 * time.Second)) {
			continue
		}

		res = append(res, instance.hostPort)
	}

	return res, nil
}
