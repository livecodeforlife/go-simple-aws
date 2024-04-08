package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/livecodeforlife/go-simple-aws/pkg/gocloud/aws/cloud"
	"github.com/livecodeforlife/go-simple-aws/pkg/gocloud/aws/provider"
	"github.com/livecodeforlife/go-simple-aws/pkg/gocloud/aws/types"
	"github.com/livecodeforlife/go-simple-aws/pkg/gocloud/core/planner"
	"github.com/livecodeforlife/go-simple-aws/pkg/gocloud/core/store"
)

func main() {
	cloud := unwrap(cloud.New(
		provider.NewResourceProvider(unwrap(config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-2")))),
		store.New(store.NewFileLoader("/tmp/myapp.json")),
		planner.NewSimplePlanner(),
	))
	myvpc := unwrap(
		cloud.CreateVPC("myvpc", &types.VpcInput{
			CidrBlock: aws.String("10.0.0.0/16"),
		}),
	)
	mysubnet := unwrap(
		cloud.CreateSubnet("mysubnet", &types.SubnetInput{
			CidrBlock: aws.String("10.0.0.0/24"),
		}),
	)
	cloud.SetSubnetVpc(mysubnet, myvpc)

	log.Println("Before Apply")
	if err := cloud.Apply(); err != nil {
		log.Fatalf("%s", err)
	}
}

func unwrap[T any](data T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return data
}
