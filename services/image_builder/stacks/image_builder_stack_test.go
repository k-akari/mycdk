package stacks_test

import (
	"testing"

	"github.com/k-akari/services/image_builder/stacks"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	"github.com/aws/jsii-runtime-go"
)

func TestImageBuilderStack(t *testing.T) {
	app := awscdk.NewApp(nil)
	stack := stacks.NewImageBuilderStack(app, "MyStack", nil)

	template := assertions.Template_FromStack(stack)

	template.ResourceCountIs(jsii.String("AWS::EC2::VPC"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::EC2::Subnet"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::EC2::InternetGateway"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::EC2::VPCGatewayAttachment"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::EC2::RouteTable"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::EC2::Route"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::EC2::SubnetRouteTableAssociation"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::CodeBuild::Project"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::ECR::Repository"), jsii.Number(1));

	template.HasResourceProperties(jsii.String("AWS::EC2::VPC"), map[string]interface{}{
		"CidrBlock": "10.0.0.0/16",
	})
	template.HasResourceProperties(jsii.String("AWS::EC2::Subnet"), map[string]interface{}{
		"AvailabilityZone": assertions.Match_AnyValue(),
		"CidrBlock": "10.0.0.0/24",
	})
	template.HasResourceProperties(jsii.String("AWS::CodeBuild::Project"), map[string]interface{}{
		"Artifacts": map[string]string{"Type": "NO_ARTIFACTS"},
		"Environment": map[string]interface{}{
			"ComputeType": "BUILD_GENERAL1_SMALL",
			"Image": "aws/codebuild/standard:1.0",
			"PrivilegedMode": true,
			"Type": "LINUX_CONTAINER",
		},
		"Source": map[string]interface{}{
			"GitCloneDepth": 1,
			"Location": "https://github.com/k-akari/ent-example.git",
			"ReportBuildStatus": true,
			"Type": "GITHUB",
		},
		"Cache": map[string]string{"Type": "NO_CACHE"},
		"Name": "ImageBuilerEntExample",
		"SourceVersion": "main",
		"Triggers": map[string]interface{}{"Webhook": true},
	})
	template.HasResourceProperties(jsii.String("AWS::ECR::Repository"), map[string]interface{}{
		"ImageScanningConfiguration": map[string]interface{}{"ScanOnPush": true},
		"LifecyclePolicy": map[string]interface{}{"LifecyclePolicyText": assertions.Match_AnyValue()},
		"RepositoryName": "ent-example",
	})
}
