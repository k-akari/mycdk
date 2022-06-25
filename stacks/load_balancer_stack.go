package stacks

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	elbv2 "github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

func NewLoadBalancerStack(scope constructs.Construct, id string, vpc ec2.Vpc, props *cdk.StackProps) (cdk.Stack, ec2.SecurityGroup) {
	var sprops cdk.StackProps
	if props != nil {
		sprops = *props
	}
	stack := cdk.NewStack(scope, &id, &sprops)

	// ALBのセキュリティグループの作成
	sg := ec2.NewSecurityGroup(stack, jsii.String("SecurityGroup"), &ec2.SecurityGroupProps{
		Vpc: vpc,
		AllowAllOutbound: jsii.Bool(true),
		Description: jsii.String("Security Group for Application Load Balancer"),
		SecurityGroupName: jsii.String("SecurityGroupForALB"),
	})
	sg.AddIngressRule(ec2.Peer_AnyIpv4(), ec2.Port_Tcp(jsii.Number(80)), jsii.String("Allow HTTP access"), jsii.Bool(false))
	sg.AddIngressRule(ec2.Peer_AnyIpv6(), ec2.Port_Tcp(jsii.Number(443)), jsii.String("Allow HTTPS access"), jsii.Bool(false))

	// ALBの作成
	elbv2.NewApplicationLoadBalancer(stack, jsii.String("ApplicationLoadBalancer"), &elbv2.ApplicationLoadBalancerProps{
		Vpc: vpc,
		DeletionProtection: jsii.Bool(false),
		InternetFacing: jsii.Bool(true),
		LoadBalancerName: jsii.String("alb-eks"),
		VpcSubnets: &ec2.SubnetSelection{
			SubnetType: ec2.SubnetType_PUBLIC,
		},
		Http2Enabled: jsii.Bool(true),
		IdleTimeout: cdk.Duration_Minutes(jsii.Number(5)),
		SecurityGroup: sg,
	})

	return stack, sg
}