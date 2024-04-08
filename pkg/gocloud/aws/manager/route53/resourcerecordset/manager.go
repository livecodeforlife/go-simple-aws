package route53resourcerecodsetmanager

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/livecodeforlife/go-simple-aws/pkg/gocloud/aws/types"
)

// NewFromConfig Creates a new instance of the resource manager
func NewFromConfig(config aws.Config) types.DnsRecordSetResourceManager {
	return &manager{
		client: route53.NewFromConfig(config),
	}
}

type manager struct {
	client *route53.Client
}

func (rm *manager) Create(input *types.DnsRecordSetInput) (*types.DnsRecordSetID, *types.DnsRecordSetOutput, error) {
	output, err := rm.client.ChangeResourceRecordSets(context.TODO(), input)
	if err != nil {
		return aws.String(""), nil, err
	}
	return output.ChangeInfo.Id, output.ChangeInfo, nil
}

func (rm *manager) Update(id *types.DnsRecordSetID, input *types.DnsRecordSetInput) (*types.DnsRecordSetID, *types.DnsRecordSetOutput, error) {
	return nil, nil, fmt.Errorf("TODO: Need to implement")
}

func (rm *manager) Retrieve(id *types.DnsRecordSetID) (*types.DnsRecordSetOutput, error) {
	output, err := rm.client.GetChange(context.TODO(), &route53.GetChangeInput{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return output.ChangeInfo, nil
}

func (rm *manager) Delete(id *types.DnsRecordSetID) (bool, error) {
	return false, fmt.Errorf("TODO: Need to implement")
}
