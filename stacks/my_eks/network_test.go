package my_eks_test

import (
	"os"
	"testing"

	myeks "mycdk/stacks/my_eks"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	assertions "github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	jsii "github.com/aws/jsii-runtime-go"
)

func TestNewNetworkStack(t *testing.T) {
	app := cdk.NewApp(nil)

	// テスト対象のスタックテンプレートを用意
	testStack := cdk.NewStack(app, jsii.String("TestStack"), &cdk.StackProps{Env: &cdk.Environment{
	 	Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	 	Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	},})
	myeks.NewNetwork(testStack)
	template := assertions.Template_FromStack(testStack)

	// 作成されるリソース数を確認
	template.ResourceCountIs(jsii.String("AWS::EC2::VPC"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::EC2::Subnet"), jsii.Number(6));
	template.ResourceCountIs(jsii.String("AWS::EC2::RouteTable"), jsii.Number(6));
	template.ResourceCountIs(jsii.String("AWS::EC2::SubnetRouteTableAssociation"), jsii.Number(6));
	template.ResourceCountIs(jsii.String("AWS::EC2::NatGateway"), jsii.Number(2));
	template.ResourceCountIs(jsii.String("AWS::EC2::EIP"), jsii.Number(2));

	// 作成されるリソースのプロパティを確認
	template.HasResourceProperties(jsii.String("AWS::EC2::VPC"), map[string]interface{}{
		"CidrBlock": "10.0.0.0/16",
		"EnableDnsHostnames": true,
		"EnableDnsSupport": true,
		"InstanceTenancy": "default",
		"Tags": []map[string]string{
			{"Key": "Name", "Value": "vpc-for-eks-cluster"},
		},
	})
}
