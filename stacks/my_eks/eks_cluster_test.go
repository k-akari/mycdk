package my_eks_test

import (
	myeks "mycdk/stacks/my_eks"
	"os"
	"testing"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	assertions "github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	jsii "github.com/aws/jsii-runtime-go"
)

func TestNewEksClusterStack(t *testing.T) {
	app := cdk.NewApp(nil)

	// テスト対象のスタックテンプレートを用意
	testStack := cdk.NewStack(app, jsii.String("TestStack"), &cdk.StackProps{Env: &cdk.Environment{
	 	Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	 	Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	},})
	vpc := ec2.NewVpc(testStack, jsii.String("VPC"), &ec2.VpcProps{
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
	myeks.NewEksCluster(testStack, vpc)
	template := assertions.Template_FromStack(testStack)

	// 作成されるリソース数を確認
	template.ResourceCountIs(jsii.String("Custom::AWSCDK-EKS-Cluster"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::EKS::Nodegroup"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::EC2::VPCEndpoint"), jsii.Number(1));

	// 作成されるリソースのプロパティを確認
	template.HasResourceProperties(jsii.String("Custom::AWSCDK-EKS-Cluster"), map[string]interface{}{
		"Config": map[string]interface{}{
			"name": "eks-cluster",
			"version": "1.22",
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
		"InstanceTypes": []string{"m5a.large", "m5.large", "m5ad.large", "m5d.large", "m5n.large", "m5dn.large"},
		"Labels": map[string]string{"app": "practice"},
		"NodegroupName": "eks-node-group",
		"ScalingConfig": map[string]float64{"DesiredSize": 2, "MaxSize": 6, "MinSize": 2,},
		"Tags": map[string]string{"Environment": "production", "Service": "service_name"},
	})
	template.HasResourceProperties(jsii.String("AWS::EC2::VPCEndpoint"), map[string]interface{}{
		"ServiceName": "com.amazonaws.ap-northeast-1.ecr.api",
		"VpcId": map[string]interface{}{"Ref": assertions.Match_StringLikeRegexp(jsii.String("VPC"))},
		"VpcEndpointType": "Interface",
	})
}
