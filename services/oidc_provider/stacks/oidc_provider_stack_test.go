package stacks_test

import (
	"testing"

	"github.com/k-akari/services/oidc_provider/stacks"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	assertions "github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	"github.com/aws/jsii-runtime-go"
)

func TestOidcProviderStack(t *testing.T) {
	app := awscdk.NewApp(nil)
	stack := stacks.NewOidcProviderStack(app, "MyStack", nil)

	template := assertions.Template_FromStack(stack)
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
