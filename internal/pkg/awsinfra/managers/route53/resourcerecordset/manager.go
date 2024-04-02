package route53resourcerecodsetmanager

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/livecodeforlife/go-simple-aws/internal/pkg/awsinfra"
)

// New Creates a new instsance of the resource manager
func New(client *route53.Client) awsinfra.ResourceManager[*route53.ChangeResourceRecordSetsInput, *types.ChangeInfo] {
	return &manager{
		client,
	}
}

type manager struct {
	client *route53.Client
}

func (rm *manager) Create(input *route53.ChangeResourceRecordSetsInput) (awsinfra.ExternalID, *types.ChangeInfo, error) {
	output, err := rm.client.ChangeResourceRecordSets(context.TODO(), input)
	if err != nil {
		return aws.String(""), nil, err
	}
	return output.ChangeInfo.Id, output.ChangeInfo, nil
}

func (rm *manager) Update(input *route53.ChangeResourceRecordSetsInput, last *types.ChangeInfo) (awsinfra.ExternalID, *types.ChangeInfo, error) {
	return aws.String(""), &types.ChangeInfo{}, fmt.Errorf("TODO: Need to implement")
}

func (rm *manager) Load(id awsinfra.ExternalID) (*types.ChangeInfo, error) {
	output, err := rm.client.GetChange(context.TODO(), &route53.GetChangeInput{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return output.ChangeInfo, nil
}

func (rm *manager) Destroy(id awsinfra.ExternalID) error {
	return fmt.Errorf("TODO: Need to implement")
}
