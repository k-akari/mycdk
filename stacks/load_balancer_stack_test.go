package stacks_test

import (
	"testing"

	"mycdk/stacks"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	assertions "github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	jsii "github.com/aws/jsii-runtime-go"
)

func TestNewLoadBalancerStack(t *testing.T) {
	app := cdk.NewApp(nil)

	// クロススタック参照のデータを用意
	refStack := cdk.NewStack(app, jsii.String("ReferenceStack"), nil)
	vpc := ec2.NewVpc(refStack, jsii.String("VPC"), &ec2.VpcProps{
		SubnetConfiguration: &[]*ec2.SubnetConfiguration{
			{
				CidrMask: jsii.Number(24),
				Name: jsii.String("subnet-for-eks-cluster-public"),
				SubnetType: ec2.SubnetType_PUBLIC,
			},
			{
				CidrMask: jsii.Number(24),
				Name: jsii.String("subnet-for-eks-cluster-private-with-nat"),
				SubnetType: ec2.SubnetType_PRIVATE_WITH_NAT,
			},
			{
				CidrMask: jsii.Number(24),
				Name: jsii.String("subnet-for-eks-cluster-private-isolated"),
				SubnetType: ec2.SubnetType_PRIVATE_ISOLATED,
			},
		},
	})

	// テスト対象のスタックとテンプレートを用意
	testStack, _ := stacks.NewLoadBalancerStack(app, "TestStack", vpc, nil)
	template := assertions.Template_FromStack(testStack)

	// 作成されるリソース数を確認
	template.ResourceCountIs(jsii.String("AWS::EC2::SecurityGroup"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::ElasticLoadBalancingV2::LoadBalancer"), jsii.Number(1));

	// 作成されるリソースのプロパティを確認
	template.HasResourceProperties(jsii.String("AWS::EC2::SecurityGroup"), map[string]interface{}{
		"GroupDescription": assertions.Match_AnyValue(),
		"GroupName": "SecurityGroupForALB",
		"SecurityGroupIngress": []map[string]interface{}{
			{
				"CidrIp": "0.0.0.0/0",
				"Description": assertions.Match_AnyValue(),
				"FromPort": 80,
				"ToPort": 80,
				"IpProtocol": "tcp",
			},
			{
				"CidrIpv6": "::/0",
				"Description": assertions.Match_AnyValue(),
				"FromPort": 443,
				"ToPort": 443,
				"IpProtocol": "tcp",
			},
		},
		"VpcId": map[string]interface{}{
			"Fn::ImportValue": assertions.Match_AnyValue(),
		},
	})
	template.HasResourceProperties(jsii.String("AWS::ElasticLoadBalancingV2::LoadBalancer"), map[string]interface{}{
		"LoadBalancerAttributes": []map[string]string{
			{"Key": "deletion_protection.enabled", "Value": "false"},
			{"Key": "idle_timeout.timeout_seconds", "Value": "300"},
		},
		"Name": "alb-eks",
		"Scheme": "internet-facing",
		"Type": "application",
	})
}
