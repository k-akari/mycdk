package stacks_test

import (
	"testing"

	"mycdk/stacks"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	assertions "github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	jsii "github.com/aws/jsii-runtime-go"
)

func TestNewEksClusterStack(t *testing.T) {
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
	sgAlb := ec2.NewSecurityGroup(refStack, jsii.String("SecurityGroupForALB"), &ec2.SecurityGroupProps{
		Vpc: vpc,
		AllowAllOutbound: jsii.Bool(true),
		Description: jsii.String("Security Group for Application Load Balancer"),
		SecurityGroupName: jsii.String("SecurityGroupForALB"),
	})

	// テスト対象のスタックとテンプレートを用意
	testStack, _ := stacks.NewEksClusterStack(app, "TestStack", vpc, sgAlb, nil)
	template := assertions.Template_FromStack(testStack)

	// 作成されるリソース数を確認
	template.ResourceCountIs(jsii.String("Custom::AWSCDK-EKS-Cluster"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::EKS::Nodegroup"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::EC2::SecurityGroupIngress"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::EC2::LaunchTemplate"), jsii.Number(1));

	// 作成されるリソースのプロパティを確認
	template.HasResourceProperties(jsii.String("Custom::AWSCDK-EKS-Cluster"), map[string]interface{}{
		"Config": map[string]interface{}{
			"name": "eks-cluster",
			"version": "1.21",
			"resourcesVpcConfig": map[string]bool{
				"endpointPublicAccess": true,
				"endpointPrivateAccess": false,
			},
		},
	})
	template.HasResourceProperties(jsii.String("AWS::EKS::Nodegroup"), map[string]interface{}{
		"ClusterName": map[string]interface{}{"Ref": assertions.Match_StringLikeRegexp(jsii.String("EKSCluster"))},
		"AmiType": "AL2_x86_64",
		"CapacityType": "SPOT",
		"ForceUpdateEnabled": true,
		"InstanceTypes": []string{"t2.micro", "t2.small", "t2.medium", "t3.micro", "t3.small", "t3.medium"},
		"Labels": map[string]string{"app": "practice"},
		"NodegroupName": "eks-node-group",
		"ScalingConfig": map[string]float64{"DesiredSize": 3, "MaxSize": 6, "MinSize": 3,},
		"Tags": map[string]string{"Environment": "production", "Service": "service_name"},
	})
	template.HasResourceProperties(jsii.String("AWS::EC2::SecurityGroupIngress"), map[string]interface{}{
		"IpProtocol": "tcp",
		"FromPort": 80,
		"SourceSecurityGroupId": map[string]interface{}{
			"Fn::ImportValue": assertions.Match_AnyValue(),
		},
		"ToPort": 80,
	})
}
