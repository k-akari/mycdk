package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func NewTestEnvStack(scope constructs.Construct, id string, props *awscdk.StackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = *props
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// Create 1 VPC with 1 public subnet
	vpc := awsec2.NewVpc(stack, jsii.String("VPCTestEnv"), &awsec2.VpcProps{
		Cidr: jsii.String("10.0.0.0/16"),
		MaxAzs: jsii.Number(1),
		SubnetConfiguration: &[]*awsec2.SubnetConfiguration{
			{
				CidrMask: jsii.Number(24),
				Name: jsii.String("subnet-for-test-env"),
				SubnetType: awsec2.SubnetType_PUBLIC,
			},
		},
		VpcName: jsii.String("vpc-for-test-env"),
	})

	// Create a role to attach an EC2 Instance
    role := awsiam.NewRole(stack, jsii.String("RoleEC2TestEnv"), &awsiam.RoleProps{
      	AssumedBy: awsiam.NewServicePrincipal(jsii.String("ec2.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
		Description: jsii.String("Iam Role for EC2 Instance"),
      	Path: jsii.String("/"),
      	RoleName: jsii.String("role-ec2-for-test-env"),
		ManagedPolicies: &[]awsiam.IManagedPolicy{
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonSSMManagedInstanceCore")),
		},
    })

	// Create an instance 
	instance := awsec2.NewInstance(stack, jsii.String("EC2InstanceTestEnv"), &awsec2.InstanceProps{
		AllowAllOutbound: jsii.Bool(true),
		InstanceType: awsec2.NewInstanceType(jsii.String("t3.small")),
		MachineImage: awsec2.NewAmazonLinuxImage(&awsec2.AmazonLinuxImageProps{
			Generation: awsec2.AmazonLinuxGeneration_AMAZON_LINUX_2,
		}),
		Role: role,
		Vpc: vpc,
		VpcSubnets: &awsec2.SubnetSelection{
			Subnets: vpc.PublicSubnets(),
		},
	})

	// Create an elastic ip and attach it to the instance
	awsec2.NewCfnEIP(stack, jsii.String("EIPEC2TestEnv"), &awsec2.CfnEIPProps{
		Domain: jsii.String("vpc"),
		InstanceId: instance.InstanceId(),
	})

	return stack
}