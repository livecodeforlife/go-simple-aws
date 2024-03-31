package main

import (
	"log"

	"github.com/livecodeforlife/go-simple-aws/internal/pkg/infra"
	awsp "github.com/livecodeforlife/go-simple-aws/internal/pkg/infra/provider/aws"
)

func main() {
	provider, err := awsp.New("us-east-2")
	if err != nil {
		log.Fatal("could not create aws provider")
	}
	myinfra := infra.New(provider)
	myvpc, err := myinfra.CreateVPC("myvpc", infra.VPC{
		CidrBlock: "10.0.0.0/16",
	})
	if err != nil {
		log.Fatal("could not create vpc")
	}
	log.Printf("%s %v", myvpc.GetID(), myvpc.GetHandler())
}
