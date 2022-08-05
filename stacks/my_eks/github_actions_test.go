package my_eks_test

import (
	myeks "mycdk/stacks/my_eks"
	"os"
	"testing"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	assertions "github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	jsii "github.com/aws/jsii-runtime-go"
)

func TestGitHubActionsStack(t *testing.T) {
	app := cdk.NewApp(nil)

	// テスト対象のスタックテンプレートを用意
	props := &cdk.StackProps{Env: &cdk.Environment{
	 	Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	 	Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	},}
	testStack := cdk.NewStack(app, jsii.String("TestStack"), props)
	myeks.NewGitHubActions(testStack)
	template := assertions.Template_FromStack(testStack)

	// 作成されるリソース数を確認

	// 作成されるリソースのプロパティを確認
	template.HasResourceProperties(jsii.String("Custom::AWSCDKOpenIdConnectProvider"), map[string]interface{}{
		"ServiceToken": assertions.Match_AnyValue(),
		"ClientIDList": []string{"sts.amazonaws.com"},
		"ThumbprintList": []string{
			"a031c46782e6e6c662c2c87c76da9aa62ccabd8e",
			"6938fd4d98bab03faadb97b34396831e3780aea1",
		},
		"Url": "https://token.actions.githubusercontent.com",
	})
	template.HasResourceProperties(jsii.String("AWS::IAM::ManagedPolicy"), map[string]interface{}{
		"PolicyDocument": map[string]interface{}{
			"Statement": assertions.Match_AnyValue(),
			"Version": assertions.Match_AnyValue(),
		},
		"Description": assertions.Match_AnyValue(),
		"ManagedPolicyName": "policy-github",
		"Path": "/",
	})
	template.HasResourceProperties(jsii.String("AWS::IAM::Role"), map[string]interface{}{
		"AssumeRolePolicyDocument": map[string]interface{}{
			"Statement": assertions.Match_AnyValue(),
			"Version": assertions.Match_AnyValue(),
		},
		"Description": assertions.Match_AnyValue(),
		"ManagedPolicyArns": assertions.Match_AnyValue(),
		"Path": "/",
		"RoleName": "role-github",
	})
}
