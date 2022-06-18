package stacks

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

func NewMyEc2Stack(scope constructs.Construct, id string, props *cdk.StackProps) cdk.Stack {
	var sprops cdk.StackProps
	if props != nil {
		sprops = *props
	}
	stack := cdk.NewStack(scope, &id, &sprops)

	// Create 1 VPC with 1 public subnet
	vpc := ec2.NewVpc(stack, jsii.String("VPCTestEnv"), &ec2.VpcProps{
		Cidr: jsii.String("10.0.0.0/16"),
		MaxAzs: jsii.Number(1),
		SubnetConfiguration: &[]*ec2.SubnetConfiguration{
			{
				CidrMask: jsii.Number(24),
				Name: jsii.String("subnet-for-test-env"),
				SubnetType: ec2.SubnetType_PUBLIC,
			},
		},
		VpcName: jsii.String("vpc-for-test-env"),
	})

	// Create a role to attach an EC2 Instance
    role := iam.NewRole(stack, jsii.String("RoleEC2TestEnv"), &iam.RoleProps{
      	AssumedBy: iam.NewServicePrincipal(jsii.String("ec2.amazonaws.com"), &iam.ServicePrincipalOpts{}),
		Description: jsii.String("Iam Role for EC2 Instance"),
      	Path: jsii.String("/"),
      	RoleName: jsii.String("role-ec2-for-test-env"),
		ManagedPolicies: &[]iam.IManagedPolicy{
			iam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonSSMManagedInstanceCore")),
		},
    })

	// Create an instance 
	instance := ec2.NewInstance(stack, jsii.String("EC2InstanceTestEnv"), &ec2.InstanceProps{
		AllowAllOutbound: jsii.Bool(true),
		InstanceType: ec2.NewInstanceType(jsii.String("t3.small")),
		MachineImage: ec2.NewAmazonLinuxImage(&ec2.AmazonLinuxImageProps{
			Generation: ec2.AmazonLinuxGeneration_AMAZON_LINUX_2,
		}),
		Role: role,
		Vpc: vpc,
		VpcSubnets: &ec2.SubnetSelection{
			Subnets: vpc.PublicSubnets(),
		},
	})

	// Create an elastic ip and attach it to the instance
	ec2.NewCfnEIP(stack, jsii.String("EIPEC2TestEnv"), &ec2.CfnEIPProps{
		Domain: jsii.String("vpc"),
		InstanceId: instance.InstanceId(),
	})

	return stack
}