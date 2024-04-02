package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	_, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-2"))
	if err != nil {
		log.Fatal("could not load aws config")
	}
	/*
		//TODO; Create a Resource Store
			myinfra := awsinfra.New(provider.NewResourceProvider(cfg))
			myvpc, err := myinfra.CreateVPC("myvpc", infra.VPC{
				CidrBlock: "10.0.0.0/16",
			})
			if err != nil {
				log.Fatal("could not create vpc")
			}
			log.Printf("%s %v", myvpc.GetID(), myvpc.GetHandler())
	*/
}
