package aws

import (
	"fmt"

	"github.com/livecodeforlife/go-simple-aws/internal/pkg/infra"
)

// New returns a new AWSProvider
func New() infra.Provider {
	return &Provider{}
}

// Provider implements the infra.Provider interface
type Provider struct {
}

// VPC returns a VPC Manager
func (p *Provider) VPC() infra.ResourceManager[infra.VPC] {
	return &vpc{}
}

// DNS returns a DNS Manager
func (p *Provider) DNS() infra.ResourceManager[infra.DNS] {
	return &dns{}
}

// Subnet returns a Subnet Manager
func (p *Provider) Subnet() infra.ResourceManager[infra.Subnet] {
	return &subnet{}
}

// AutoScale returns a Subnet Manager
func (p *Provider) AutoScale() infra.ResourceManager[infra.AutoScale] {
	return &autoscale{}
}

// LoadBalancer returns a Subnet Manager
func (p *Provider) LoadBalancer() infra.ResourceManager[infra.LoadBalancer] {
	return &loadBalancer{}
}

// LaunchTemplate returns a Subnet Manager
func (p *Provider) LaunchTemplate() infra.ResourceManager[infra.LaunchTemplate] {
	return &launchTemplate{}
}

type vpc struct{}

func (rm *vpc) Create(id string, r infra.VPC) (*infra.Resource[infra.VPC], error) {
	return nil, fmt.Errorf("TODO: Need to implement")
}

func (rm *vpc) Delete(id string) error {
	return fmt.Errorf("TODO: Need to implement")
}

type dns struct{}

func (rm *dns) Create(id string, r infra.DNS) (*infra.Resource[infra.DNS], error) {
	return nil, fmt.Errorf("TODO: Need to implement")
}

func (rm *dns) Delete(id string) error {
	return fmt.Errorf("TODO: Need to implement")
}

type subnet struct{}

func (rm *subnet) Create(id string, r infra.Subnet) (*infra.Resource[infra.Subnet], error) {
	return nil, fmt.Errorf("TODO: Need to implement")
}

// Delete simulates a resource deletion
func (rm *subnet) Delete(id string) error {
	return fmt.Errorf("TODO: Need to implement")
}

type autoscale struct{}

func (rm *autoscale) Create(id string, r infra.AutoScale) (*infra.Resource[infra.AutoScale], error) {
	return nil, fmt.Errorf("TODO: Need to implement")
}

// Delete simulates a resource deletion
func (rm *autoscale) Delete(id string) error {
	return fmt.Errorf("TODO: Need to implement")
}

type loadBalancer struct{}

func (rm *loadBalancer) Create(id string, r infra.LoadBalancer) (*infra.Resource[infra.LoadBalancer], error) {
	return nil, fmt.Errorf("TODO: Need to implement")
}

// Delete simulates a resource deletion
func (rm *loadBalancer) Delete(id string) error {
	return fmt.Errorf("TODO: Need to implement")
}

type launchTemplate struct{}

func (rm *launchTemplate) Create(id string, r infra.LaunchTemplate) (*infra.Resource[infra.LaunchTemplate], error) {
	return nil, fmt.Errorf("TODO: Need to implement")
}

// Delete simulates a resource deletion
func (rm *launchTemplate) Delete(id string) error {
	return fmt.Errorf("TODO: Need to implement")
}
