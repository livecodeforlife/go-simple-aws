# go-simple-aws


Pseudo Code:

# Assignment

Please design a program to do a zero downtime deploy using standard EC2 and related resources such as AMI, ELB, etc. using an AWS SDK or similar package. Please do not use Terraform, Cloudformation, Kubernetes or other such tool.

# Requirements:

- Pseudo code for the full program.
- Full directory structure as if this was indeed a full implementation as a part of a larger application to perform many infrastructure related tasks.
- Discussion of approach.
- A full implementation of one or two interesting parts of the program including tests such as those which require mocking out multiple external calls.


# Pseudo code for full program

var Dns
var Network
var SecurityGroup
var NetworkACL
var IAMPolicy
var MyServiceElb
var MyServiceElbIamPolice
var MyServiceElbIamRole
var MyServiceELBGreen
var MyServiceElbGreenPolice
var MyServiceElbGreenRole
var MyServiceELBBlue
var MyServiceElbBluePolice
var MyServiceElbBlueRole
var MyServiceInstanceLaunchTemplate
var MyServiceAutoScalingGroup
var MyServiceAmi
var MyServiceAmiSecurityGroup
var MyServiceAmiIamPolice
var MyServiceAmiIamRole
var MyServiceInstance
var MyServiceInstanceSecurityGroup
var MyServiceInstanceIamPolice
var MyServiceInstanceIamRole


Pseudo code:

Ensure main VPC is created
    Ensure VPC Security Group is created
    Ensure VPC Network ACL is created
Ensure DNS is created
    Ensure DNS certificate is created
Ensure the main Load Balancer is created
    Ensure public subnet for main Load Balancer is created
        Ensure Public Subnet Security Group is created
        Ensure Public Subnet Network Acl is created
    Ensure DNS is pointing to the main Load Balancer


let mainActualTarget = getMainActualTarget(MyServiceElb)
if mainActualTarget = green 
    then 
        deploy(MyServiceElbBlue)
        setMainActualTarget(MyServiceElb) = blue
    else
        deploy(MyServiceElbGreen)
        setMainActualTarget(MyServiceElb) = green

func deploy(MyServiceElbBlue, imageIdHash)
    Ensure private subnet





    Set Launch Template with imageHash
    Set AutoScalingBlue
    Ensure Launch Security Group is created

        




