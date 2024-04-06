package main

import (
	"log"

	"github.com/livecodeforlife/go-simple-aws/pkg/awsinfra"
	"github.com/livecodeforlife/go-simple-aws/pkg/coreinfra"
)

func main() {
	//TODO; Create a Resource Store
	infra := awsinfra.NewWithDefaults()
	myvpc, err := infra.CreateVPC("myvpc", &awsinfra.CreateVpcInput{
		CidrBlock: awsinfra.String("10.0.0.0/16"),
	})
	if err != nil {
		log.Fatal("could not create vpc")
	}
	mysubnet, err := infra.CreateSubnet("mysubnet", &awsinfra.CreateSubnetInput{})
	if err != nil {
		log.Fatal("could not create vpc")
	}
	coreinfra.AddDependency(
		infra.ResourceStorer(),
		mysubnet,
		myvpc,
		func(subnet *awsinfra.CreateSubnetInput, vpc *awsinfra.VpcResource) error {
			subnet.VpcId = vpc.Output().VpcId
			return nil
		})
	infra.Apply()
	infra.Destroy()
}
