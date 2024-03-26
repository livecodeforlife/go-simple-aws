package infra

// VPCDefinition has the configuration to create a VPC
type VPCDefinition struct {
}

// VPCHandle is the object that the VPC Creator returns when it creates a VPC
type VPCHandle struct {
}

// VPCCreator is the interface that does create the VPC
type VPCCreator interface {
	create(vpc VPCDefinition) (*VPCHandle, error)
}

func createOrUpdateVpc(def VPCDefinition, vpc VPCCreator) (*VPCHandle, error) {
	return vpc.create(def)
}
