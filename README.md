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

OK -> let mainVpc = createOrUpdateVpc(VPCDefinition)
let mainDns = createOrUpdateDns(DNSDefinition)
let mainSubnet = createOrUpdateSubnet(SubnetDefinition(vpc))
let mainElb = createOrUpdateElb(ELBDefinition(mainDns,mainSubnet,mainVpc))

let greenSubnet = createOrUpdateSubnet(SubnetDefinition(vpc))
let greenElb = createOrUpdateElb(ElbDefinition(greenSubnet,mainVpc))
let greenLaunchTemplate = createOrUpdateASGLaunchTemplate(ElbLaunchTemplate)
let greenAutoScale = createOrUpdateGreenAutoScale(greenAutoScale,greenElb)

let blueSubnet = createOrUpdateSubnet(SubnetDefinition(vpc))
let blueElb = createOrUpdateElb(ElbDefinition(blueSubnet,mainVpc))
let blueLaunchTemplate = createOrUpdateASGLaunchTemplate(ElbLaunchTemplate)
let blueAutoScale = createOrUpdateGreenAutoScale(greenAutoScale,blueElb)

let image = compile image from source code
let deployLauchTemplate;
let deployElb;
if mainElb.getActiveEnvironment is "blue" {
    deployLaunchTemplate = greenLaunchTemplate
    deployElb = greenElb
    deployAutoScale = greenAutoScale
} else {
    deployLaunchTemplate = blueLaunchTemplate
    deployElb = blueElb
    deployAutoScale = blueAutoScale
}
perform_smoke_tests(deployElb)
perform_switch(mainElb,deployElb)


createOrUpdateVpc(VPCDefinition):

createOrUpdateDns(DNSDefinition):

createOrUpdateSubnet(SubnetDefinition(vpc)):

createOrUpdateElb(ELBDefinition(mainDns,mainSubnet,mainVpc)):

createOrUpdateSubnet(SubnetDefinition(vpc)):

createOrUpdateElb(ElbDefinition(blueSubnet,mainVpc)):

createOrUpdateASGLaunchTemplate(ElbLaunchTemplate):

createOrUpdateGreenAutoScale(greenAutoScale,blueElb):