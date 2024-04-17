package discovery

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Registry interface {
	//register instance
	Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error
	//deregister instance
	Deregister(ctx context.Context, instanceID string, serviceName string) error
	//get addresses for all service instances
	GetServiceAddresses(ctx context.Context, serviceID string) ([]string, error)
	//notify registry of instance healthy state
	ReportHealthyState(instanceID string, serviceName string) error
}

var ErrNotFound = errors.New("service addresses not found")

func GenerateInstanceId(serviceName string) string {
	integer := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	return fmt.Sprintf("%s-%d", serviceName, integer)
}
