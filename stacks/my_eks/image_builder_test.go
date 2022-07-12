package my_eks_test

import (
	myeks "mycdk/stacks/my_eks"
	"os"
	"testing"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	jsii "github.com/aws/jsii-runtime-go"
)

func TestImageBuilderStack(t *testing.T) {
	app := cdk.NewApp(nil)
	
	// テスト対象のスタックテンプレートを用意
	props := &cdk.StackProps{Env: &cdk.Environment{
	 	Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	 	Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	},}
	testStack := cdk.NewStack(app, jsii.String("TestStack"), props)
	myeks.NewImageBuilder(testStack, props)
	template := assertions.Template_FromStack(testStack)

	// 作成されるリソース数を確認
	template.ResourceCountIs(jsii.String("AWS::IAM::Role"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::CodeBuild::Project"), jsii.Number(1));
	template.ResourceCountIs(jsii.String("AWS::ECR::Repository"), jsii.Number(3));

	// 作成されるリソースのプロパティを確認
	template.HasResourceProperties(jsii.String("AWS::IAM::Role"), map[string]interface{}{
		"AssumeRolePolicyDocument": map[string]interface{}{
			"Statement": []map[string]interface{}{
				{
					"Action": "sts:AssumeRole",
					"Effect": "Allow",
					"Principal": map[string]string{
						"Service": "codebuild.amazonaws.com",
					},
				},
			},
			"Version": assertions.Match_AnyValue(),
		},
		"Description": assertions.Match_AnyValue(),
		"Path": "/",
		"RoleName": "role-codebuild-for-image-builder",
	})
	template.HasResourceProperties(jsii.String("AWS::CodeBuild::Project"), map[string]interface{}{
		"Artifacts": map[string]string{"Type": "NO_ARTIFACTS"},
		"Environment": map[string]interface{}{
			"ComputeType": "BUILD_GENERAL1_SMALL",
			"EnvironmentVariables": []map[string]string{
				{"Name": "AWS_ACCOUNT", "Type": "PLAINTEXT", "Value": os.Getenv("CDK_DEFAULT_ACCOUNT")},
				{"Name": "AWS_REGION", "Type": "PLAINTEXT", "Value": os.Getenv("CDK_DEFAULT_REGION")},
			},
			"Image": "aws/codebuild/standard:1.0",
			"PrivilegedMode": true,
			"Type": "LINUX_CONTAINER",
		},
		"Source": map[string]interface{}{
			"BuildSpec": "buildspec.yml",
			"GitCloneDepth": 1,
			"Location": "https://github.com/k-akari/ent-example.git",
			"ReportBuildStatus": true,
			"Type": "GITHUB",
		},
		"Cache": map[string]string{"Type": "NO_CACHE"},
		"ConcurrentBuildLimit": 1,
		"Name": "EKSImageBuildProject",
		"QueuedTimeoutInMinutes": 60,
		"SourceVersion": "main",
		"TimeoutInMinutes": 20,
		"Triggers": map[string]interface{}{"Webhook": true},
	})
	template.HasResourceProperties(jsii.String("AWS::ECR::Repository"), map[string]interface{}{
		"ImageScanningConfiguration": map[string]interface{}{"ScanOnPush": true},
		"LifecyclePolicy": map[string]interface{}{"LifecyclePolicyText": assertions.Match_AnyValue()},
		"RepositoryName": "eks-app",
	})
	template.HasResourceProperties(jsii.String("AWS::ECR::Repository"), map[string]interface{}{
		"ImageScanningConfiguration": map[string]interface{}{"ScanOnPush": true},
		"LifecyclePolicy": map[string]interface{}{"LifecyclePolicyText": assertions.Match_AnyValue()},
		"RepositoryName": "eks-migration",
	})
	template.HasResourceProperties(jsii.String("AWS::ECR::Repository"), map[string]interface{}{
		"ImageScanningConfiguration": map[string]interface{}{"ScanOnPush": true},
		"LifecyclePolicy": map[string]interface{}{"LifecyclePolicyText": assertions.Match_AnyValue()},
		"RepositoryName": "eks-web",
	})
}
