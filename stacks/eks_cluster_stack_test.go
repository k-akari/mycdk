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
	networkStack := cdk.NewStack(app, jsii.String("NetworkStack"), nil)
	vpc := ec2.NewVpc(networkStack, jsii.String("VPC"), &ec2.VpcProps{
		Cidr: jsii.String("10.0.0.0/16"),
		MaxAzs: jsii.Number(3),
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
		VpcName: jsii.String("vpc-for-eks-cluster"),
	})

	eksClusterStack, _ := stacks.NewEksClusterStack(app, "MyStack", vpc, nil)

	template := assertions.Template_FromStack(eksClusterStack)

	template.ResourceCountIs(jsii.String("Custom::AWSCDK-EKS-Cluster"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::EKS::Nodegroup"), jsii.Number(1));

	template.HasResourceProperties(jsii.String("AWS::EKS::Nodegroup"), map[string]interface{}{
		"ClusterName": map[string]interface{}{"Ref": assertions.Match_StringLikeRegexp(jsii.String("EKSCluster"))},
		"AmiType": "AL2_x86_64",
		"CapacityType": "SPOT",
		"DiskSize": 10,
		"ForceUpdateEnabled": true,
	})
}
