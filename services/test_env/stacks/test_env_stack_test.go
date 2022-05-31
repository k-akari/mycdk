package stacks_test

import (
	"testing"

	"github.com/k-akari/services/test_env/stacks"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	assertions "github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	"github.com/aws/jsii-runtime-go"
)

func TestTestEnvStack(t *testing.T) {
	app := awscdk.NewApp(nil)
	stack := stacks.NewTestEnvStack(app, "MyStack", nil)

	template := assertions.Template_FromStack(stack)

	template.ResourceCountIs(jsii.String("AWS::EC2::VPC"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::EC2::Subnet"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::EC2::InternetGateway"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::EC2::VPCGatewayAttachment"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::EC2::Route"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::EC2::SubnetRouteTableAssociation"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::EC2::Instance"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::IAM::Role"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::IAM::InstanceProfile"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::EC2::SecurityGroup"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::EC2::EIP"), jsii.Number(1));

	template.HasResourceProperties(jsii.String("AWS::EC2::VPC"), map[string]interface{}{
		"CidrBlock": "10.0.0.0/16",
	})
	template.HasResourceProperties(jsii.String("AWS::EC2::Subnet"), map[string]interface{}{
		"AvailabilityZone": assertions.Match_AnyValue(),
		"CidrBlock": "10.0.0.0/24",
	})
	template.HasResourceProperties(jsii.String("AWS::EC2::Instance"), map[string]interface{}{
		"InstanceType": "t3.small",
	})
	template.HasResourceProperties(jsii.String("AWS::EC2::EIP"), map[string]interface{}{
		"Domain": "vpc",
		"InstanceId": assertions.Match_AnyValue(),
	})
}
