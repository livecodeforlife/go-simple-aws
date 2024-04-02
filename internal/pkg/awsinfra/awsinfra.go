package awsinfra

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	autoscalingtypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
)

// Infra manages the lifecycle of cloud infrastructure components, such as creation,
// update, and deletion.
type Infra struct {
	defaultRollback  bool                      //Experimental
	resourceProvider ResourceProvider          // Provider implements resource creation interfaces.
	resourceStore    ResourceStore             //Permanent datastore the syncs the infra
	localStore       map[InternalID]ExternalID // Tracks created resources and avoid duplicated resources
	resourceStack    resourceStack             //remembers the sequence of created elements to rollback
}

// New initializes a new infrastructure manager with the specified provider.
func New(resourceProvicer ResourceProvider, resourceStore ResourceStore, withRollback bool) *Infra {
	return &Infra{
		withRollback,
		resourceProvicer,
		resourceStore,
		make(map[InternalID]ExternalID),
		resourceStack{},
	}
}

// ExternalID is a unique identifier for a resource in an external system.
type ExternalID = *string

// InternalID is a unique identifier for a resource in the infrastructure manager.
type InternalID = string

// ResourceManager create or update resources
type ResourceManager[Input any, Output any] interface {
	Create(input Input) (ExternalID, Output, error)
	Update(input Input, last Output) (ExternalID, Output, error)
	Load(id ExternalID) (Output, error)
	ResourceDestroyer
}

// ResourceDestroyer destroy resources
type ResourceDestroyer interface {
	Destroy(id ExternalID) error
}

// ResourceStore helps with idempotency
type ResourceStore interface {
	Exists(internalID InternalID) (bool, error)
	Get(internalID InternalID) (ExternalID, error)
	Set(internalID InternalID, externalID ExternalID) error
}

// ResourceProvider aggregates interfaces for creating cloud resources. Implementations of ResourceProvider
// enable the creation of VPCs, DNS records, and subnets, along with managing their resource handlers.
type ResourceProvider interface {
	VPC() ResourceManager[*ec2.CreateVpcInput, *ec2types.Vpc]
	DNSRecordSet() ResourceManager[*route53.ChangeResourceRecordSetsInput, *route53types.ChangeInfo]
	Subnet() ResourceManager[*ec2.CreateSubnetInput, *ec2types.Subnet]
	LoadBalancer() ResourceManager[*elbv2.CreateLoadBalancerInput, []elbv2types.LoadBalancer]
	LaunchTemplate() ResourceManager[*ec2.CreateLaunchTemplateInput, *ec2types.LaunchTemplate]
	AutoScalingGroup() ResourceManager[*autoscaling.CreateAutoScalingGroupInput, *autoscalingtypes.AutoScalingGroup]
}

// CreateVPC requests the creation of a VPC resource in the cloud, using the provided definition.
func (i *Infra) CreateVPC(id string, input *ec2.CreateVpcInput) (*ec2types.Vpc, error) {
	return createWithRollback(i, id, input, i.resourceProvider.VPC())
}

// CreateDNS requests the creation of a DNS record in the cloud, using the provided definition.
func (i *Infra) CreateDNS(id string, input *route53.ChangeResourceRecordSetsInput) (*route53types.ChangeInfo, error) {
	return createWithRollback(i, id, input, i.resourceProvider.DNSRecordSet())
}

// CreateSubnet requests the creation of a Subnet resource in the cloud, using the provided definition.
func (i *Infra) CreateSubnet(id string, input *ec2.CreateSubnetInput) (*ec2types.Subnet, error) {
	return createWithRollback(i, id, input, i.resourceProvider.Subnet())
}

// CreateLoadBalancer requests the creation of a Subnet resource in the cloud, using the provided definition.
func (i *Infra) CreateLoadBalancer(id string, input *elbv2.CreateLoadBalancerInput) ([]elbv2types.LoadBalancer, error) {
	return createWithRollback(i, id, input, i.resourceProvider.LoadBalancer())
}

// CreateLaunchTemplate requests the creation of a LaunchTemplate resource in the cloud, using the provided definition.
func (i *Infra) CreateLaunchTemplate(id string, input *ec2.CreateLaunchTemplateInput) (*ec2types.LaunchTemplate, error) {
	return createWithRollback(i, id, input, i.resourceProvider.LaunchTemplate())
}

// CreateAutoScale requests the creation of a LaunchTemplate resource in the cloud, using the provided definition.
func (i *Infra) CreateAutoScale(id string, input *autoscaling.CreateAutoScalingGroupInput) (*autoscalingtypes.AutoScalingGroup, error) {
	return createWithRollback(i, id, input, i.resourceProvider.AutoScalingGroup())
}

func (i *Infra) validateID(id string) error {
	if id == "" {
		return &InfraError{ErrBlankResourceID, nil}
	}
	if _, exists := i.localStore[id]; exists {
		return &InfraError{ErrResourceExists, nil}
	}
	return nil
}

func (i *Infra) validateInitialization() error {
	if i.resourceProvider == nil {
		return &InfraError{ErrMissingResourceProvider, nil}
	}
	if i.localStore == nil {
		return &InfraError{ErrMissingLocalStore, nil}
	}
	if i.resourceStore == nil {
		return &InfraError{ErrMissingResourceStore, nil}
	}
	return nil
}

// Destroy all elements under the resource stack
func (i *Infra) Destroy() error {
	for {
		rs, err := i.resourceStack.Pop()
		if err != nil {
			break
		}
		//Destroy the cloud resource. it takes the externalID from the localStore
		if err := rs.resourceDestroyer.Destroy(i.localStore[rs.id]); err != nil {
			return &InfraError{ErrFailedResourceManagerDestroy, fmt.Errorf("ID: %s, Caused by %v ", rs.id, err)}
		}
		delete(i.localStore, rs.id) //deletes the id from the localStore, so it can be reused
	}
	return nil
}

func createWithRollback[Input any, Output any](infra *Infra, id InternalID, input Input, resourceManager ResourceManager[Input, Output]) (Output, error) {
	output, err := create(infra, id, input, resourceManager)
	if err != nil {
		if !infra.defaultRollback {
			return output, err
		}
		if err := infra.Destroy(); err != nil { //Destroy all stacked resources
			return output, err
		}
	}
	return output, err
}

// create is a generic function that encapsulates common logic for resource creation.
// It checks for the presence of a provider and resources, ensuring id uniqueness and provider
// ability to innerCreate and store the resource.
func create[Input any, Output any](infra *Infra, id InternalID, input Input, resourceManager ResourceManager[Input, Output]) (Output, error) {
	var output Output
	var outputID ExternalID

	if err := infra.validateInitialization(); err != nil {
		return output, err
	}
	if err := infra.validateID((id)); err != nil {
		return output, err
	}
	exists, err := infra.resourceStore.Exists(id)
	if err != nil {
		return output, &InfraError{ErrFailedResourceStoreExists, fmt.Errorf("ID: %s, Caused by %v ", id, err)}
	}
	if !exists {
		//Creates the resource
		externalID, created, err := resourceManager.Create(input)
		if err != nil {
			return output, &InfraError{ErrFailedResourceManagerCreate, fmt.Errorf("ID: %s, Caused by %v ", id, err)}
		}
		//Pushes the resource to the resource stack, allowing for rollback in case of error
		infra.resourceStack.Push(&ResourceState{
			resourceDestroyer: resourceManager,
			id:                id,
		})
		//Set the externalID to the external resourceStore
		if err := infra.resourceStore.Set(id, externalID); err != nil {
			return output, &InfraError{ErrFailedResourceStoreSet, fmt.Errorf("ID: %s, Caused by %v ", id, err)}
		}
		//Updates the return values
		output = created
		outputID = externalID
	} else {
		//Get the last created resource ID from the external resource store
		lastID, err := infra.resourceStore.Get(id)
		if err != nil {
			return output, &InfraError{ErrFailedResourceStoreGet, fmt.Errorf("ID: %s, Caused by %v ", id, err)}
		}
		//Loads the resource using the last external ID
		last, err := resourceManager.Load(lastID)
		if err != nil {
			return output, &InfraError{ErrFailedResourceManagerLoad, fmt.Errorf("ID: %s, Caused by %v ", id, err)}
		}
		//Updates the resource, merging the input with last element
		//Sometimes update is not possible, then a deletion and creation may happen
		//In that case externalID may change
		externalID, updated, err := resourceManager.Update(input, last)
		if err != nil {
			return output, &InfraError{ErrFailedResourceManagerUpdate, fmt.Errorf("ID: %s, Caused by %v ", id, err)}
		}
		//Pushes the resource to the resource stack, allowing for rollback in case of error
		infra.resourceStack.Push(&ResourceState{
			resourceDestroyer: resourceManager,
			id:                id,
		})
		//Set the externalID to the external resourceStore only if it has changed
		if *externalID != *lastID {
			if err := infra.resourceStore.Set(id, externalID); err != nil {
				return output, &InfraError{ErrFailedResourceStoreSet, fmt.Errorf("ID: %s, Caused by %v ", id, err)}
			}
		}
		output = updated
		outputID = externalID
	}
	//Store the new id into the localStore. this will prevent new creations with the same ID in the same execution
	infra.localStore[id] = outputID
	return output, nil
}

// ResourceState represents the resource that was already processed
type ResourceState struct {
	resourceDestroyer ResourceDestroyer
	id                InternalID
}
type resourceStack []*ResourceState

// Push adds an element to the top of the stack.
func (s *resourceStack) Push(val *ResourceState) {
	*s = append(*s, val)
}

// Pop removes and returns the top element from the stack.
func (s *resourceStack) Pop() (*ResourceState, error) {
	if len(*s) == 0 {
		return nil, fmt.Errorf("stack is empty")
	}
	lastIndex := len(*s) - 1
	val := (*s)[lastIndex]
	*s = (*s)[:lastIndex]
	return val, nil
}

// InfraError is the error generated by the infra package
type InfraError struct {
	Code     int
	CausedBy error
}

// Error returns the error message associated with the error code.
func (e InfraError) Error() string {
	switch e.Code {
	case ErrMissingResourceProvider:
		return fmt.Sprintf("Provider is missing")
	case ErrMissingLocalStore:
		return "LocalStore is missing"
	case ErrMissingResourceStore:
		return "ResourceStore is missing"
	case ErrResourceExists:
		return "There is already another resource with the same id"
	case ErrBlankResourceID:
		return "The resource id is blank"
	case ErrFailedResourceManagerDestroy:
		return fmt.Sprintf("Failed to destroy reource; %s", e.CausedBy)
	case ErrFailedResourceManagerLoad:
		return fmt.Sprintf("Failed to load reource; %s", e.CausedBy)
	case ErrFailedResourceManagerUpdate:
		return fmt.Sprintf("Failed to update reource; %s", e.CausedBy)
	case ErrFailedResourceManagerCreate:
		return fmt.Sprintf("Failed to create reource; %s", e.CausedBy)
	case ErrFailedResourceStoreSet:
		return fmt.Sprintf("Failed to set resource into the store; %s", e.CausedBy)
	case ErrFailedResourceStoreGet:
		return fmt.Sprintf("Failed to get resource from the store; %s", e.CausedBy)
	case ErrFailedResourceStoreExists:
		return fmt.Sprintf("Failed to check if resource exists in the store; %s", e.CausedBy)
	default:
		return "Unknown error"
	}
}

const (
	//ErrMissingResourceProvider is the error code for missing resource provider
	ErrMissingResourceProvider = iota
	//ErrMissingLocalStore is the error code for missing local store
	ErrMissingLocalStore
	//ErrMissingResourceStore is the error code for missing resource store
	ErrMissingResourceStore
	//ErrResourceExists is the error code for existing resource
	ErrResourceExists
	//ErrBlankResourceID is the error code for blank resource id
	ErrBlankResourceID
	//ErrFailedResourceManagerDestroy is the error code for failed resource manager destroy
	ErrFailedResourceManagerDestroy
	//ErrFailedResourceManagerLoad is the error code for failed resource manager load
	ErrFailedResourceManagerLoad
	//ErrFailedResourceManagerUpdate is the error code for failed resource manager update
	ErrFailedResourceManagerUpdate
	//ErrFailedResourceManagerCreate is the error code for failed resource manager create
	ErrFailedResourceManagerCreate
	//ErrFailedResourceStoreSet is the error code for failed resource store set
	ErrFailedResourceStoreSet
	//ErrFailedResourceStoreGet is the error code for failed resource store get
	ErrFailedResourceStoreGet
	//ErrFailedResourceStoreExists is the error code for failed resource store exists
	ErrFailedResourceStoreExists
)
