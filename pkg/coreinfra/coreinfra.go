package coreinfra

import (
	"encoding/json"
	"fmt"
)

// AddDependency creates a dependency to a Resource Input before its creation
func AddDependency[MyInput, Input any, Output any, EID any](storer ResourceStorer, self *LazyResource[MyInput], other *LazyResource[Input], applyF func(MyInput, *Resource[Input, Output, EID]) error) {
	self.dependencies = append(
		self.dependencies,
		dependency{
			id: other.ID(),
			applyF: func() error {
				if applyF != nil {
					otherResource, err := getStoreResource[Input, Output, EID](storer, other.ID())
					if err != nil {
						return err
					}
					if err := applyF(self.input, otherResource); err != nil {
						return err
					}

				}
				return nil
			}})
}

// CreateResource is a function that creates a future resource
func CreateResource[Input any, Output any, EID any](planner Planner, storer ResourceStorer, manager ResourceManager[Input, Output, EID], id ID, input Input, dependsOn ...ID) (*LazyResource[Input], error) {
	if planner == nil {
		return nil, &Error{ErrMissingResourcePlanner, nil}
	}
	if storer == nil {
		return nil, &Error{ErrMissingResourceStore, nil}
	}
	if manager == nil {
		return nil, &Error{ErrMissingResourceManager, nil}
	}
	if id == "" {
		return nil, &Error{ErrBlankResourceID, nil}
	}
	resource := &LazyResource[Input]{
		id:           id,
		dependencies: []dependency{},
		createFn: func() error {
			_, err := createOrUpdateResourceStrict[Input, Output, EID](storer, manager, id, input)
			return err
		},
		deleteFn: func() (bool, error) {
			return deleteResourceStrict[Input, Output, EID](storer, manager, id)
		},
	}
	if err := planner.AddResource(resource); err != nil {
		return nil, &Error{ErrResourceCreateNotAuthorized, fmt.Errorf("ID: %s, Caused by %v ", id, err)}
	}

	return resource, nil
}

// ApplyPlan build all planned resources
func ApplyPlan(planner Planner) error {
	for _, resource := range planner.TopoSortForCreation() {
		if err := resource.getCreateFn()(); err != nil {
			return err
		}
	}
	return nil
}

// DestroyPlan all planned resources
func DestroyPlan(planner Planner) error {
	for _, resource := range planner.TopoSortForDeletion() {
		if _, err := resource.getDeleteFn(); err != nil {
			return err
		}
	}
	return nil
}

// NewSimplePlanner creates a new planner to be used inside the infra package
func NewSimplePlanner() Planner {
	return &simplePlanner{}
}

type simplePlanner struct {
	resources []lazyResource
}

func (p *simplePlanner) AddResource(resource lazyResource) error {
	p.resources = append(p.resources, resource)
	return nil
}

func (p *simplePlanner) TopoSortForCreation() []lazyResource {
	return p.resources
}

func (p *simplePlanner) TopoSortForDeletion() []lazyResource {
	return reverseSliceCopy(p.resources)
}

func reverseSliceCopy[T any](s []T) []T {
	sc := make([]T, len(s))
	copy(sc, s)
	for i, j := 0, len(sc)-1; i < j; i, j = i+1, j-1 {
		sc[i], sc[j] = sc[j], sc[i]
	}
	return sc
}

// LazyResource is a resource to be generated
type LazyResource[Input any] struct {
	id           ID
	input        Input
	dependencies []dependency
	createFn     func() error
	deleteFn     func() (bool, error)
}

type dependency struct {
	id     ID
	applyF func() error
}

// ID returns the resource ID
func (r *LazyResource[Input]) ID() ID {
	return r.id
}
func (r *LazyResource[Input]) getCreateFn() func() error {
	return r.createFn
}
func (r *LazyResource[Input]) getDeleteFn() (bool, error) {
	return r.deleteFn()
}
func (r *LazyResource[Input]) getDependencies() []dependency {
	return r.dependencies
}

type lazyResource interface {
	ID() ID
	getDependencies() []dependency
	getCreateFn() func() error
	getDeleteFn() (bool, error)
}

// ID represents a unique identifier within this package.
type ID string

// ExtID is an interface for managing external IDs given by cloud providers,
// allowing conversion between a generic type and its string representation.
type ExtID[EID any] interface {
	FromStr(string) EID
	ToStr(EID) string
}

// ResourceManager is the interface that outlines methods for creating, retrieving,
// updating, and deleting resources, abstracting over specific cloud provider implementations.
type ResourceManager[Input any, Output any, EID any] interface {
	Create(Input) (ExtID[EID], Output, error)
	Retrieve(ExtID[EID]) (Output, error)
	Update(ExtID[EID], Input) (ExtID[EID], Output, error)
	Delete(ExtID[EID]) (bool, error)
}

// ResourceStorer is the interface for persisting resource information, providing
// methods for checking existence, retrieving, storing, and deleting resource data.
type ResourceStorer interface {
	Exists(ID) (bool, error)
	Get(ID) ([]byte, error)
	Set(ID, []byte) error
	Delete(ID) error
}

// Planner is a planner that creates elements in the right order, respecting dependencies
type Planner interface {
	AddResource(lazyResource) error
	TopoSortForCreation() []lazyResource
	TopoSortForDeletion() []lazyResource
}

// Resource represents an initialized cloud resource, encapsulating its input and output data,
// unique identifiers, and dependencies on other resources.
type Resource[Input any, Output any, EID any] struct {
	id        ID
	extID     ExtID[EID]
	input     Input
	output    Output
	dependsOn []ID
}

// createOrUpdateResourceStrict encapsulates the logic for creating or updating a resource.
// It performs authorization checks, ensures uniqueness, and delegates to the ResourceManager
// for creation or updating actions.
// Checks for the presence of necessary components and validates the provided ID.
// On passing the checks, it either creates a new resource or updates an existing one,
// finally persisting the resource state.
func createOrUpdateResourceStrict[Input any, Output any, EID any](storer ResourceStorer, manager ResourceManager[Input, Output, EID], id ID, input Input) (*Resource[Input, Output, EID], error) {

	if storer == nil {
		return nil, &Error{ErrMissingResourceStore, nil}
	}

	if manager == nil {
		return nil, &Error{ErrMissingResourceManager, nil}
	}

	if id == "" {
		return nil, &Error{ErrBlankResourceID, nil}
	}

	var resource *Resource[Input, Output, EID]

	exists, err := storer.Exists(id)
	if err != nil {
		return nil, &Error{ErrResourceStoreExists, fmt.Errorf("ID: %s, Caused by %v ", id, err)}
	}

	if !exists {
		resource, err = createResource(manager, id, input)
		if err != nil {
			return nil, err
		}
	} else {
		resource, err = updateResource[Input, Output, EID](storer, manager, id, input)
		if err != nil {
			return nil, err
		}
	}
	if err = setStoreResource(storer, resource); err != nil {
		return nil, err
	}
	return resource, nil
}

// deleteResourceStrict encapsulates the logic for safely deleting a resource.
// It checks for authorization, verifies the resource's existence, and if authorized,
// proceeds with deletion using the ResourceManager.
func deleteResourceStrict[Input any, Output any, EID any](storer ResourceStorer, manager ResourceManager[Input, Output, EID], id ID) (bool, error) {

	// Implementation ensures the presence of necessary components and valid ID.
	// If checks pass, it proceeds to delete the resource and notifies the authorizer.
	if storer == nil {
		return false, &Error{ErrMissingResourceStore, nil}
	}
	if manager == nil {
		return false, &Error{ErrMissingResourceManager, nil}
	}
	if id == "" {
		return false, &Error{ErrBlankResourceID, nil}
	}
	exists, err := storer.Exists(id)
	if err != nil {
		return false, &Error{ErrResourceStoreExists, fmt.Errorf("ID: %s, Caused by %v ", id, err)}
	}
	if !exists {
		//return false because did not update
		return false, nil
	}
	deleted, err := deleteResource[Input, Output, EID](storer, manager, id)
	if err != nil {
		return false, err
	}
	return deleted, nil

}

// ID returns the unique identifier of the resource.
func (r *Resource[Input, Output, EID]) ID() ID {
	return r.id
}

// ExtID returns the external ID of the resource provided by the cloud provider.
func (r *Resource[Input, Output, EID]) ExtID() ExtID[EID] {
	return r.extID
}

// Input returns the input data used to create or update the resource.
func (r *Resource[Input, Output, EID]) Input() Input {
	return r.input
}

// Output returns the output data from the cloud provider about the resource.
func (r *Resource[Input, Output, EID]) Output() Output {
	return r.output
}

// DependsOn returns a list of IDs that the resource depends on.
func (r *Resource[Input, Output, EID]) DependsOn() []ID {
	return r.dependsOn
}

// ToJSON serializes the Resource into JSON bytes for storage.
func (r *Resource[Input, Output, EID]) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

// FromJSON deserializes JSON bytes back into a Resource struct.
func (r *Resource[Input, Output, EID]) FromJSON(b []byte) error {
	return json.Unmarshal(b, r)
}

func deleteResource[Input any, Output any, EID any](storer ResourceStorer, manager ResourceManager[Input, Output, EID], id ID) (bool, error) {
	resource, err := getStoreResource[Input, Output, EID](storer, id)
	if err != nil {
		return false, err
	}
	deleted, err := manager.Delete(resource.extID)
	if err != nil {
		return false, &Error{ErrResourceManagerUpdate, fmt.Errorf("ID: %s, Caused by %v ", id, err)}
	}

	return deleted, nil
}

func updateResource[Input any, Output any, EID any](storer ResourceStorer, manager ResourceManager[Input, Output, EID], id ID, input Input) (*Resource[Input, Output, EID], error) {
	resource, err := getStoreResource[Input, Output, EID](storer, id)
	if err != nil {
		return nil, err
	}
	extID, output, err := manager.Update(resource.extID, input)
	if err != nil {
		return nil, &Error{ErrResourceManagerUpdate, fmt.Errorf("ID: %s, Caused by %v ", id, err)}
	}
	return &Resource[Input, Output, EID]{
		id:        id,
		extID:     extID,
		input:     input,
		output:    output,
		dependsOn: resource.dependsOn,
	}, nil
}

func setStoreResource[Input any, Output any, EID any](storer ResourceStorer, resource *Resource[Input, Output, EID]) error {
	resourceAsBytes, err := resource.ToJSON()
	if err != nil {
		return &Error{ErrResourceToBytes, fmt.Errorf("ID: %s, Caused by %v ", resource.id, err)}
	}
	if err = storer.Set(resource.id, resourceAsBytes); err != nil {
		return &Error{ErrResourceStoreSet, fmt.Errorf("ID: %s, Caused by %v ", resource.id, err)}
	}
	return nil
}

func deleteStoreResource[Input any, Output any, EID any](storer ResourceStorer, resource *Resource[Input, Output, EID]) error {
	err := storer.Delete(resource.id)
	if err != nil {
		return &Error{ErrResourceStoreDelete, fmt.Errorf("ID: %s, Caused by %v ", resource.id, err)}
	}
	return nil
}

func getStoreResource[Input any, Output any, EID any](storer ResourceStorer, id ID) (*Resource[Input, Output, EID], error) {
	var resource *Resource[Input, Output, EID]
	resourceJSON, err := storer.Get(id)
	if err != nil {
		return nil, &Error{ErrResourceStoreGet, fmt.Errorf("ID: %s, Caused by %v ", id, err)}
	}
	if err := resource.FromJSON(resourceJSON); err != nil {
		return nil, &Error{ErrResourceLoadFromBytes, fmt.Errorf("ID: %s, Caused by %v ", id, err)}
	}
	return resource, nil
}

func createResource[Input any, Output any, EID any](manager ResourceManager[Input, Output, EID], id ID, input Input) (*Resource[Input, Output, EID], error) {
	extID, output, err := manager.Create(input)
	if err != nil {
		return nil, &Error{ErrResourceManagerCreate, fmt.Errorf("ID: %s, Caused by %v ", id, err)}
	}
	return &Resource[Input, Output, EID]{
		id:     id,
		extID:  extID,
		input:  input,
		output: output,
	}, nil
}

// Error represents an error within the infra package, encapsulating an error code
// and an optional underlying cause for more detailed error handling.
type Error struct {
	Code     int   // The error code representing the specific error condition.
	CausedBy error // An optional underlying error that caused this error, if any.
}

// Error returns a human-readable message associated with the error code,
// optionally including the message from an underlying error, if present.
func (e Error) Error() string {
	baseMessage := "Infra error: " // Prefix for all error messages to identify the source.

	// Switch on the error code to provide a specific error message.
	switch e.Code {
	case ErrMissingResourceManager:
		return baseMessage + "Resource Manager is missing."
	case ErrMissingResourcePlanner:
		return baseMessage + "Resource Authorizer is missing."
	case ErrMissingResourceStore:
		return baseMessage + "Resource Store is missing."
	case ErrResourceCreateNotAuthorized:
		return baseMessage + fmt.Sprintf("Creation not authorized: %v", e.CausedBy)
	case ErrResourceDeleteNotAuthorized:
		return baseMessage + fmt.Sprintf("Deletion not authorized: %v", e.CausedBy)
	case ErrBlankResourceID:
		return baseMessage + "Resource ID is blank."
	case ErrResourceManagerDestroy:
		return baseMessage + fmt.Sprintf("Failed to destroy resource: %v", e.CausedBy)
	case ErrResourceManagerLoad:
		return baseMessage + fmt.Sprintf("Failed to load resource: %v", e.CausedBy)
	case ErrResourceManagerUpdate:
		return baseMessage + fmt.Sprintf("Failed to update resource: %v", e.CausedBy)
	case ErrResourceManagerCreate:
		return baseMessage + fmt.Sprintf("Failed to create resource: %v", e.CausedBy)
	case ErrResourceStoreSet:
		return baseMessage + fmt.Sprintf("Failed to store resource: %v", e.CausedBy)
	case ErrResourceStoreDelete:
		return baseMessage + fmt.Sprintf("Failed to delete resource from store: %v", e.CausedBy)
	case ErrResourceStoreGet:
		return baseMessage + fmt.Sprintf("Failed to retrieve resource from store: %v", e.CausedBy)
	case ErrResourceStoreExists:
		return baseMessage + fmt.Sprintf("Failed to check resource existence in store: %v", e.CausedBy)
	case ErrResourceLoadFromBytes:
		return baseMessage + fmt.Sprintf("Failed to deserialize resource from bytes: %v", e.CausedBy)
	case ErrResourceToBytes:
		return baseMessage + fmt.Sprintf("Failed to serialize resource to bytes: %v", e.CausedBy)
	default:
		return baseMessage + "Unknown error."
	}
}

const (
	// ErrMissingResourceManager indicates a missing Resource Manager.
	ErrMissingResourceManager = iota

	// ErrMissingResourcePlanner indicates a missing Resource Authorizer.
	ErrMissingResourcePlanner

	// ErrMissingResourceStore indicates a missing Resource Store.
	ErrMissingResourceStore

	// ErrResourceCreateNotAuthorized indicates creation is not authorized.
	ErrResourceCreateNotAuthorized

	// ErrResourceDeleteNotAuthorized indicates deletion is not authorized.
	ErrResourceDeleteNotAuthorized

	// ErrBlankResourceID indicates the Resource ID is blank.
	ErrBlankResourceID

	// ErrResourceManagerDestroy indicates failure in destroying a resource.
	ErrResourceManagerDestroy

	// ErrResourceManagerLoad indicates failure in loading a resource.
	ErrResourceManagerLoad

	// ErrResourceManagerUpdate indicates failure in updating a resource.
	ErrResourceManagerUpdate

	// ErrResourceManagerCreate indicates failure in creating a resource.
	ErrResourceManagerCreate

	// ErrResourceStoreSet indicates failure in storing a resource.
	ErrResourceStoreSet

	// ErrResourceStoreDelete indicates failure in deleting a resource from the store.
	ErrResourceStoreDelete

	// ErrResourceStoreGet indicates failure in retrieving a resource from the store.
	ErrResourceStoreGet

	// ErrResourceStoreExists indicates failure in checking if a resource exists in the store.
	ErrResourceStoreExists

	// ErrResourceLoadFromBytes indicates failure in deserializing a resource from bytes.
	ErrResourceLoadFromBytes

	// ErrResourceToBytes indicates failure in serializing a resource to bytes.
	ErrResourceToBytes
)
