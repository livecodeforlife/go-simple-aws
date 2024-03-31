package awsinfra

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	autoscalingtypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/stretchr/testify/assert"
)

type TResourceStore struct {
	existsErr error
	getErr    error
	setErr    error
	store     map[InternalID]ExternalID
}

func (rs *TResourceStore) Exists(internalID InternalID) (bool, error) {
	_, ok := rs.store[internalID]
	return ok, rs.existsErr
}
func (rs *TResourceStore) Get(internalID InternalID) (ExternalID, error) {
	return rs.store[internalID], rs.getErr
}
func (rs *TResourceStore) Set(internalID InternalID, externalID ExternalID) error {
	rs.store[internalID] = externalID
	return rs.setErr
}

// TResourceManager mocks the creation of VPC resources for testing.
type TResourceManager[I any, O any] struct {
	Output     O
	CreateErr  error
	LoadErr    error
	UpdateErr  error
	DestroyErr error
	Eid        ExternalID
	creates    uint
	updates    uint
	loads      uint
	deletes    uint
}

func (rm *TResourceManager[Input, Output]) Create(input Input) (ExternalID, Output, error) {
	rm.creates++
	return rm.Eid, rm.Output, rm.CreateErr
}

func (rm *TResourceManager[Input, Output]) Update(input Input, last Output) (ExternalID, Output, error) {
	rm.updates++
	return rm.Eid, rm.Output, rm.UpdateErr
}
func (rm *TResourceManager[Input, Output]) Load(id ExternalID) (Output, error) {
	rm.loads++
	return rm.Output, rm.LoadErr
}

// Destroy simulates a resource deletion
func (rm *TResourceManager[Input, Output]) Destroy(id ExternalID) error {
	return rm.DestroyErr
}

// TestProvider aggregates mocks for various resource creators and a resource store.
type TestProvider struct {
	vpc            TResourceManager[ec2.CreateVpcInput, *ec2types.Vpc]
	dns            TResourceManager[route53.ChangeResourceRecordSetsInput, *route53types.ResourceRecordSet]
	subnet         TResourceManager[ec2.CreateSubnetInput, *ec2types.Subnet]
	loadBalancer   TResourceManager[elbv2.CreateLoadBalancerInput, *elbv2types.LoadBalancer]
	launchTemplate TResourceManager[ec2.CreateLaunchTemplateInput, *ec2types.LaunchTemplate]
	autoScale      TResourceManager[autoscaling.CreateAutoScalingGroupInput, *autoscalingtypes.AutoScalingGroup]
}

func (p *TestProvider) VPC() ResourceManager[ec2.CreateVpcInput, *ec2types.Vpc] {
	return &p.vpc
}
func (p *TestProvider) DNSRecordSet() ResourceManager[route53.ChangeResourceRecordSetsInput, *route53types.ResourceRecordSet] {
	return &p.dns
}
func (p *TestProvider) Subnet() ResourceManager[ec2.CreateSubnetInput, *ec2types.Subnet] {
	return &p.subnet
}
func (p *TestProvider) LoadBalancer() ResourceManager[elbv2.CreateLoadBalancerInput, *elbv2types.LoadBalancer] {
	return &p.loadBalancer
}
func (p *TestProvider) LaunchTemplate() ResourceManager[ec2.CreateLaunchTemplateInput, *ec2types.LaunchTemplate] {
	return &p.launchTemplate
}
func (p *TestProvider) AutoScalingGroup() ResourceManager[autoscaling.CreateAutoScalingGroupInput, *autoscalingtypes.AutoScalingGroup] {
	return &p.autoScale
}

// TestCreate verifies the functionality of creating VPC, Subnet, and DNS resources
// using a mock infrastructure provider. It checks the correctness of the created
// resources' IDs and handlers against expected values.
func TestResourceCreation(t *testing.T) {
	// Setup: Initializes a new infrastructure with a mock provider containing predefined responses.
	const (
		VPCID            = "vpcid"
		DNSID            = "dnsid"
		SUBNETID         = "subnetid"
		LBID             = "ldbid"
		LAUNCHTEMPLATEID = "launchtemplateid"
		AUTOSCALEID      = "autoscaleid"
	)
	expectedStore := map[string]*string{
		VPCID:            aws.String("vpceid"),
		DNSID:            aws.String("dnseid"),
		SUBNETID:         aws.String("subneteid"),
		LBID:             aws.String("ldbeid"),
		LAUNCHTEMPLATEID: aws.String("launchtemplateeid"),
		AUTOSCALEID:      aws.String("autoscaleeid"),
	}
	eid := func(id InternalID) ExternalID {
		return expectedStore[id]
	}
	provider := &TestProvider{
		vpc:            TResourceManager[ec2.CreateVpcInput, *ec2types.Vpc]{Output: &ec2types.Vpc{}, Eid: eid(VPCID)},
		dns:            TResourceManager[route53.ChangeResourceRecordSetsInput, *route53types.ResourceRecordSet]{Output: &route53types.ResourceRecordSet{}, Eid: eid(DNSID)},
		subnet:         TResourceManager[ec2.CreateSubnetInput, *ec2types.Subnet]{Output: &ec2types.Subnet{}, Eid: eid(SUBNETID)},
		loadBalancer:   TResourceManager[elbv2.CreateLoadBalancerInput, *elbv2types.LoadBalancer]{Output: &elbv2types.LoadBalancer{}, Eid: eid(LBID)},
		launchTemplate: TResourceManager[ec2.CreateLaunchTemplateInput, *ec2types.LaunchTemplate]{Output: &ec2types.LaunchTemplate{}, Eid: eid(LAUNCHTEMPLATEID)},
		autoScale:      TResourceManager[autoscaling.CreateAutoScalingGroupInput, *autoscalingtypes.AutoScalingGroup]{Output: &autoscalingtypes.AutoScalingGroup{}, Eid: eid(AUTOSCALEID)},
	}
	store := &TResourceStore{
		store: make(map[InternalID]ExternalID),
	}
	infra := New(provider, store, false)
	testCreate(t, store, VPCID, eid(VPCID), &ec2types.Vpc{}, ec2.CreateVpcInput{}, infra.CreateVPC)
	testCreate(t, store, DNSID, eid(DNSID), &route53types.ResourceRecordSet{}, route53.ChangeResourceRecordSetsInput{}, infra.CreateDNS)
	testCreate(t, store, SUBNETID, eid(SUBNETID), &ec2types.Subnet{}, ec2.CreateSubnetInput{}, infra.CreateSubnet)
	testCreate(t, store, LBID, eid(LBID), &elbv2types.LoadBalancer{}, elbv2.CreateLoadBalancerInput{}, infra.CreateLoadBalancer)
	testCreate(t, store, LAUNCHTEMPLATEID, eid(LAUNCHTEMPLATEID), &ec2types.LaunchTemplate{}, ec2.CreateLaunchTemplateInput{}, infra.CreateLaunchTemplate)
	testCreate(t, store, AUTOSCALEID, eid(AUTOSCALEID), &autoscalingtypes.AutoScalingGroup{}, autoscaling.CreateAutoScalingGroupInput{}, infra.CreateAutoScale)
}

func testCreate[Input any, Output any](t *testing.T, store ResourceStore, id InternalID, expectedExternalID ExternalID, expectedOutput Output, input Input, create func(id InternalID, input Input) (Output, error)) {
	output, err := create(id, input)
	externalID, _ := store.Get(id)
	assert.Equal(t, externalID, expectedExternalID, "ID should match the expected value.")
	assert.Equal(t, output, expectedOutput, "output should match the expected value.")
	assert.Nil(t, err, "No error should be returned during resource creation.")
}

func Test_create(t *testing.T) {
	type args struct {
		infra           *Infra
		id              string
		input           string
		resourceManager ResourceManager[string, string]
	}
	tests := []struct {
		name        string
		args        args
		want        string
		wantErr     bool
		wantErrCode uint
	}{
		{
			name: "Given a well initialized Infra, an id and an input, should return the expected output",
			args: args{
				infra: &Infra{
					withRollback:     false,
					resourceProvider: &TestProvider{},
					resourceStore:    &TResourceStore{store: make(map[InternalID]ExternalID)},
					localStore:       make(map[string]*string),
				},
				id:    "testInternalID",
				input: "testInput",
				resourceManager: &TResourceManager[string, string]{
					Output: "testOutput",
					Eid:    aws.String("testExternalID"),
				},
			},
			want:    "testOutput",
			wantErr: false,
		},
		{
			name: "When resouce provider is null, should return an error",
			args: args{
				infra: &Infra{
					withRollback:     false,
					resourceProvider: nil,
					resourceStore:    &TResourceStore{store: make(map[InternalID]ExternalID)},
					localStore:       make(map[string]*string),
				},
				id:    "testInternalID",
				input: "testInput",
				resourceManager: &TResourceManager[string, string]{
					Output: "testOutput",
					Eid:    aws.String("testExternalID"),
				},
			},
			want:        "",
			wantErr:     true,
			wantErrCode: ErrMissingResourceProvider,
		},
		{
			name: "When resource store is null, should return an error",
			args: args{
				infra: &Infra{
					withRollback:     false,
					resourceProvider: &TestProvider{},
					resourceStore:    nil,
					localStore:       make(map[string]*string),
				},
				id:    "testInternalID",
				input: "testInput",
				resourceManager: &TResourceManager[string, string]{
					Output: "testOutput",
					Eid:    aws.String("testExternalID"),
				},
			},
			want:        "",
			wantErr:     true,
			wantErrCode: ErrMissingResourceStore,
		},
		{
			name: "When local store is nill, should return an error",
			args: args{
				infra: &Infra{
					withRollback:     false,
					resourceProvider: &TestProvider{},
					resourceStore:    &TResourceStore{store: make(map[InternalID]ExternalID)},
					localStore:       nil,
				},
				id:    "testInternalID",
				input: "testInput",
				resourceManager: &TResourceManager[string, string]{
					Output: "testOutput",
					Eid:    aws.String("testExternalID"),
				},
			},
			want:        "",
			wantErr:     true,
			wantErrCode: ErrMissingLocalStore,
		},
		{
			name: "Given a well initialized Infra, a blank id and an input, should return an error",
			args: args{
				infra: &Infra{
					withRollback:     false,
					resourceProvider: &TestProvider{},
					resourceStore:    &TResourceStore{store: make(map[InternalID]ExternalID)},
					localStore:       make(map[string]*string),
					resourceStack:    resourceStack{},
				},
				id:    "",
				input: "testInput",
				resourceManager: &TResourceManager[string, string]{
					Output: "testOutput",
					Eid:    aws.String("testExternalID"),
				},
			},
			want:        "",
			wantErr:     true,
			wantErrCode: ErrBlankResourceID,
		},
		{
			name: "When trying to create with an already existing id on local store, should generate an error",
			args: args{
				infra: &Infra{
					withRollback:     false,
					resourceProvider: &TestProvider{},
					resourceStore:    &TResourceStore{store: make(map[InternalID]ExternalID)},
					localStore: map[InternalID]ExternalID{
						"testInternalID": aws.String("testExternalID"), //id is already in the local store
					},
					resourceStack: resourceStack{},
				},
				id:    "testInternalID",
				input: "testInput",
				resourceManager: &TResourceManager[string, string]{
					Output: "testOutput",
					Eid:    aws.String("testExternalID"),
				},
			},
			want:        "",
			wantErr:     true,
			wantErrCode: ErrResourceExists,
		},
		{
			name: "When resource store generates an error, then should return an error",
			args: args{
				infra: &Infra{
					withRollback:     false,
					resourceProvider: &TestProvider{},
					resourceStore: &TResourceStore{
						existsErr: fmt.Errorf("Exists error"),
						store:     make(map[InternalID]ExternalID),
					},
					localStore:    make(map[string]*string),
					resourceStack: resourceStack{},
				},
				id:    "testInternalID",
				input: "testInput",
				resourceManager: &TResourceManager[string, string]{
					Output: "testOutput",
					Eid:    aws.String("testExternalID"),
				},
			},
			want:        "",
			wantErr:     true,
			wantErrCode: ErrFailedResourceStoreExists,
		},
		{
			name: "When resource manager fails to create, then should return an error",
			args: args{
				infra: &Infra{
					withRollback:     false,
					resourceProvider: &TestProvider{},
					resourceStore: &TResourceStore{
						existsErr: nil,
						store:     make(map[InternalID]ExternalID),
					},
					localStore:    make(map[string]*string),
					resourceStack: resourceStack{},
				},
				id:    "testInternalID",
				input: "testInput",
				resourceManager: &TResourceManager[string, string]{
					Output:    "testOutput",
					CreateErr: fmt.Errorf("Something bad has happened"),
					Eid:       aws.String("testExternalID"),
				},
			},
			want:        "",
			wantErr:     true,
			wantErrCode: ErrFailedResourceManagerCreate,
		},
		{
			name: "When resource store fails to set, then should return an error",
			args: args{
				infra: &Infra{
					withRollback:     false,
					resourceProvider: &TestProvider{},
					resourceStore: &TResourceStore{
						existsErr: nil,
						setErr:    fmt.Errorf("Something bad has happened"),
						store:     make(map[InternalID]ExternalID),
					},
					localStore:    make(map[string]*string),
					resourceStack: resourceStack{},
				},
				id:    "testInternalID",
				input: "testInput",
				resourceManager: &TResourceManager[string, string]{
					Output:    "testOutput",
					CreateErr: nil,
					Eid:       aws.String("testExternalID"),
				},
			},
			want:        "",
			wantErr:     true,
			wantErrCode: ErrFailedResourceStoreSet,
		},
		{
			name: "When resource exists on external store, if resource store fails to get, then should return an error",
			args: args{
				infra: &Infra{
					withRollback:     false,
					resourceProvider: &TestProvider{},
					resourceStore: &TResourceStore{
						existsErr: nil,
						setErr:    nil,
						getErr:    fmt.Errorf("Something bad has happened"),
						store: map[InternalID]ExternalID{
							"testInternalID": aws.String("testExternalID"), //id is already in the external store
						},
					},
					localStore:    make(map[string]*string),
					resourceStack: resourceStack{},
				},
				id:    "testInternalID",
				input: "testInput",
				resourceManager: &TResourceManager[string, string]{
					Output:    "testOutput",
					CreateErr: nil,
					Eid:       aws.String("testExternalID"),
				},
			},
			want:        "",
			wantErr:     true,
			wantErrCode: ErrFailedResourceStoreGet,
		},
		{
			name: "When resource exists on external store, if resource manager fails to Load, then should return an error",
			args: args{
				infra: &Infra{
					withRollback:     false,
					resourceProvider: &TestProvider{},
					resourceStore: &TResourceStore{
						existsErr: nil,
						setErr:    nil,
						getErr:    nil,
						store: map[InternalID]ExternalID{
							"testInternalID": aws.String("testExternalID"), //id is already in the external store
						},
					},
					localStore:    make(map[string]*string),
					resourceStack: resourceStack{},
				},
				id:    "testInternalID",
				input: "testInput",
				resourceManager: &TResourceManager[string, string]{
					Output:    "testOutput",
					CreateErr: nil,
					LoadErr:   fmt.Errorf("Something bad has happened"),
					Eid:       aws.String("testExternalID"),
				},
			},
			want:        "",
			wantErr:     true,
			wantErrCode: ErrFailedResourceManagerLoad,
		},
		{
			name: "When resource exists on external store, if resource manager fails to Update, then should return an error",
			args: args{
				infra: &Infra{
					withRollback:     false,
					resourceProvider: &TestProvider{},
					resourceStore: &TResourceStore{
						existsErr: nil,
						setErr:    nil,
						getErr:    nil,
						store: map[InternalID]ExternalID{
							"testInternalID": aws.String("testExternalID"), //id is already in the external store
						},
					},
					localStore:    make(map[string]*string),
					resourceStack: resourceStack{},
				},
				id:    "testInternalID",
				input: "testInput",
				resourceManager: &TResourceManager[string, string]{
					Output:    "testOutput",
					CreateErr: nil,
					UpdateErr: fmt.Errorf("Something bad has happened"),
					Eid:       aws.String("testExternalID"),
				},
			},
			want:        "",
			wantErr:     true,
			wantErrCode: ErrFailedResourceManagerUpdate,
		},
		{
			name: "When resource exists on external store, if resource manager returns a different ExternalID and store fails to set, then should return an error",
			args: args{
				infra: &Infra{
					withRollback:     false,
					resourceProvider: &TestProvider{},
					resourceStore: &TResourceStore{
						existsErr: nil,
						setErr:    fmt.Errorf("Something bad has happened"),
						getErr:    nil,
						store: map[InternalID]ExternalID{
							"testInternalID": aws.String("testExternalID"), //id is already in the external store
						},
					},
					localStore:    make(map[string]*string),
					resourceStack: resourceStack{},
				},
				id:    "testInternalID",
				input: "testInput",
				resourceManager: &TResourceManager[string, string]{
					Output:    "testOutput",
					CreateErr: nil,
					UpdateErr: nil,
					Eid:       aws.String("differentExternalID"),
				},
			},
			want:        "",
			wantErr:     true,
			wantErrCode: ErrFailedResourceStoreSet,
		},
		{
			name: "When resource exists on external store, if resource manager returns a different ExternalID, then should return the expected output",
			args: args{
				infra: &Infra{
					withRollback:     false,
					resourceProvider: &TestProvider{},
					resourceStore: &TResourceStore{
						existsErr: nil,
						setErr:    nil,
						getErr:    nil,
						store: map[InternalID]ExternalID{
							"testInternalID": aws.String("testExternalID"), //id is already in the external store
						},
					},
					localStore:    make(map[string]*string),
					resourceStack: resourceStack{},
				},
				id:    "testInternalID",
				input: "testInput",
				resourceManager: &TResourceManager[string, string]{
					Output:    "testOutput",
					CreateErr: nil,
					UpdateErr: nil,
					Eid:       aws.String("differentExternalID"),
				},
			},
			want:    "testOutput",
			wantErr: false,
		},
		{
			name: "When resource exists on external store, and external IDS are the same, then should return the expected output",
			args: args{
				infra: &Infra{
					withRollback:     false,
					resourceProvider: &TestProvider{},
					resourceStore: &TResourceStore{
						existsErr: nil,
						setErr:    nil,
						getErr:    nil,
						store: map[InternalID]ExternalID{
							"testInternalID": aws.String("testExternalID"), //id is already in the external store
						},
					},
					localStore:    make(map[string]*string),
					resourceStack: resourceStack{},
				},
				id:    "testInternalID",
				input: "testInput",
				resourceManager: &TResourceManager[string, string]{
					Output:    "testOutput",
					CreateErr: nil,
					UpdateErr: nil,
					Eid:       aws.String("testExternalID"),
				},
			},
			want:    "testOutput",
			wantErr: false,
		},
		{
			name: "When rollback is active, if fail after the creation, then should rollback",
			args: args{
				infra: &Infra{
					withRollback:     true,
					resourceProvider: &TestProvider{},
					resourceStore: &TResourceStore{
						existsErr: nil,
						setErr:    fmt.Errorf("Something very bad has happened"),
						getErr:    nil,
						store:     make(map[string]*string),
					},
					localStore:    make(map[string]*string),
					resourceStack: resourceStack{},
				},
				id:    "testInternalID",
				input: "testInput",
				resourceManager: &TResourceManager[string, string]{
					Output:    "testOutput",
					CreateErr: nil,
					UpdateErr: nil,
					Eid:       aws.String("testExternalID"),
				},
			},
			want:        "",
			wantErr:     true,
			wantErrCode: ErrFailedResourceStoreSet,
		},
		{
			name: "When rolling back, if destroy generated an error, then should return the error",
			args: args{
				infra: &Infra{
					withRollback:     true,
					resourceProvider: &TestProvider{},
					resourceStore: &TResourceStore{
						existsErr: nil,
						setErr:    fmt.Errorf("Something very bad has happened"),
						getErr:    nil,
						store:     make(map[string]*string),
					},
					localStore:    make(map[string]*string),
					resourceStack: resourceStack{},
				},
				id:    "testInternalID",
				input: "testInput",
				resourceManager: &TResourceManager[string, string]{
					Output:     "testOutput",
					CreateErr:  nil,
					UpdateErr:  nil,
					DestroyErr: fmt.Errorf("Somenthing bad has happened"),
					Eid:        aws.String("testExternalID"),
				},
			},
			want:        "",
			wantErr:     true,
			wantErrCode: ErrFailedResourceManagerDestroy,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createWithRollback(tt.args.infra, tt.args.id, tt.args.input, tt.args.resourceManager)
			if (err != nil) != tt.wantErr {
				t.Errorf("create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("create() = %v, want %v", got, tt.want)
			}
			if err != nil {
				switch err.(type) {
				case *InfraError:
					if err.(*InfraError).Code != int(tt.wantErrCode) {
						t.Errorf("create() error code = %v, want %v", err.(*InfraError).Code, tt.wantErrCode)
					}
					t.Logf("Well captured error %s", err)
				default:
					t.Errorf("Error is not a InfraError")
				}
			}
		})
	}
}

func TestDefaultErrorMsg(t *testing.T) {
	err := &InfraError{
		Code:     32187128709,
		CausedBy: nil,
	}
	assert.Equal(t, "Unknown error", err.Error())
}
